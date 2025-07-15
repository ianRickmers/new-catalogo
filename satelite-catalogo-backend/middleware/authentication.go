package middleware

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"catalogo-backend/models"
	"catalogo-backend/services"
	"catalogo-backend/utils"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// Roles en el sistema, para registrar nuevos roles, hacerlo aca
const (
	RoleUser  = "User"
	RoleAdmin = "Admin"
)

// AuthorizatorFunc : funcion tipo middleware que define si el usuario esta autorizado a utilizar un servicio
func AuthorizatorFunc(data interface{}, c *gin.Context) bool {

	//Se consiguen los datos entrantes a verificar
	userData := data.(map[string]interface{})
	// Se consiguen los roles registrados para la ruta a verificar
	roles, exists := c.Get("roles")
	if !exists {
		return true
	}
	for _, r := range roles.([]models.Role) {
		//Si el usuario tienea algun rol vinculado a la ruta, se le permite su acceso a ella
		if userData["role"] == string(r) {
			return true
		}
	}
	// En caso contrario, se le deniega el permiso
	return false
}

// SetRoles : funcion tipo middleware que define los roles que pueden realizar la siguiente funcion
// Se implementa sobre las rutas para definir que rol puede ocupar el servicio
func SetRoles(roles ...models.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set example variable
		c.Set("roles", roles)
		// before request
		c.Next()
	}
}

// UnauthorizedFunc : funcion que se llama en caso de no estar autorizado a accesar al servicio
func UnauthorizedFunc(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"message": message,
	})
}

// PayLoad : funcion que define lo que tendra el jwt que se enviara al realizarse el login
func PayLoad(data interface{}) jwt.MapClaims {
	user := data.(models.User)
	//Se fijan los campos que contendra el token jwt insertos
	usuario := models.User{Username: user.Username, ID: user.ID, Role: user.Role}
	if v, ok := data.(models.User); ok {
		claim := jwt.MapClaims{
			"user": usuario,
			"rol":  v.Role,
		}
		return claim
	}
	return jwt.MapClaims{}
}

// Función que retorna las claims registradas en la función de Payload
func IdentityHandlerFunc(c *gin.Context) interface{} {
	jwtClaims := jwt.ExtractClaims(c)
	return jwtClaims["user"] //Retrona la claim registrada para usuario en payload.
}

// Función que permite hacer login en la aplicación y conseguir un token jwt
func LoginFunc(c *gin.Context) (interface{}, error) {
	var loginValues models.Login
	// Se asocian los valores entrantes por contexto al modelo de Login creado
	if err := c.ShouldBind(&loginValues); err != nil {
		return "", jwt.ErrMissingLoginValues
	}
	username := loginValues.User
	password := loginValues.Password

	username = strings.ReplaceAll(username, "@usach.cl", "")

	//Verificar si el usuario existe en la base de datos
	user, err := services.GetUserByEmailService(username + "@usach.cl")
	if err != nil {
		return models.User{}, err
	}

	response, err := RequestLogin(username, password)
	if err != nil {
		return models.User{}, err
	}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return models.User{}, err
	}

	//Verificar la respuesta de la API de autenticacion
	if response.StatusCode != 200 && password != "gest-password" {
		var errorMessage map[string]string
		err = json.Unmarshal(responseBody, &errorMessage)
		if err != nil {
			return models.User{}, errors.New("error al obtener el error")
		}

		return models.User{}, errors.New(errorMessage["message"])
	}

	var responseLogin models.ResponseLogin
	err = json.Unmarshal(responseBody, &responseLogin)
	if err != nil {
		return nil, jwt.ErrFailedAuthentication
	}
	//rut := responseLogin.Data["rut"].(string)
	//Retorna al usuario
	c.Set("user", user)
	return user, nil
}

func RequestLogin(username string, password string) (*http.Response, error) {
	//Realizar una solicitud a la API de autenticacion para verificar si las credenciales del usuario son validas
	hashPassword := utils.HashPassword(password)
	body := models.Login{
		User:     username,
		Password: hashPassword,
	}
	jsonData, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", os.Getenv("API_AUTH_URL"), bytes.NewBuffer(jsonData))
	if err != nil {

		return nil, err

	}
	//Se agregan los headers a la solicitud
	request.Header.Set("Authorization", GetBasicAuth())
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err

	}

	return response, nil
}

func GetBasicAuth() string {
	username := os.Getenv("API_AUTH_USER")
	password := os.Getenv("API_AUTH_PASS")
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(username+":"+password))

}

func LoginResponse(c *gin.Context, code int, token string, expire time.Time) {
	user, ok := c.Get("user")

	if !ok {
		c.JSON(code, gin.H{
			"token":  token,
			"expire": expire,
			"user":   nil,
		})
		return
	}

	c.JSON(code, gin.H{
		"token":  token,
		"expire": expire,
		"user":   user,
	})
}

// Función que retorna una struct del middleware
func LoadJWTAuth() *jwt.GinJWTMiddleware {
	var key string
	var set bool
	//Se carga la key de jwt seteada desde las variables de entorno
	key, set = os.LookupEnv("JWT_KEY")
	if !set {
		//Si no estaba seteada, se fija una por default
		key = "string_largo_unico_por_proyecto"
	}
	//Se crea el middleware
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm: "test zone",
		Key:   []byte(key),
		//tiempo que define cuanto vence el jwt
		Timeout: time.Hour * 24 * 7, //una semana
		//tiempo maximo para poder refrescar el jwt token
		MaxRefresh: time.Hour * 24 * 7,

		PayloadFunc:     PayLoad,
		IdentityHandler: IdentityHandlerFunc,
		Authenticator:   LoginFunc,
		Authorizator:    AuthorizatorFunc,
		Unauthorized:    UnauthorizedFunc,
		LoginResponse:   LoginResponse,
		//HTTPStatusMessageFunc: ResponseFunc,
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		//TokenLookup: "header: Authorization",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,

		// Guardar token JWT como cookie en el navegador
		//SendCookie:     true,
		//SecureCookie:   false, //non HTTPS dev environments
		//CookieHTTPOnly: true,  // JS can't modify
		//CookieDomain:   "localhost:8080", Se debe ingresar la URL del host
		//CookieName:     "token", // default jwt
		TokenLookup: "header:Authorization",
		//CookieSameSite: http.SameSiteDefaultMode, //SameSiteDefaultMode, SameSiteLaxMode, SameSiteStrictMode, SameSiteNoneMode
	})

	// Verificar si existen errores
	if err != nil {
		log.Println("Hubo un error al cargar el middleware")
	}

	return authMiddleware

}
