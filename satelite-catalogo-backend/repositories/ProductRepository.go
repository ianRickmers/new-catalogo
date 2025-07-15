package repositories

import (
	"catalogo-backend/database"
	"catalogo-backend/models"
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var productRepo *ProductRepository

type ProductRepository struct {
	collection *mongo.Collection
}

func NewProductRepository() *ProductRepository {
	if database.Client == nil {
		log.Fatal("MongoDB client not initialized. Call InitMongo() first.")
	}
	if productRepo == nil {
		log.Println("Initializing ProductRepository")
		db := database.GetDatabase()
		collection := db.Collection("products")
		productRepo = &ProductRepository{collection: collection}
	}
	return productRepo
}

func (r *ProductRepository) Create(ctx context.Context, product models.Product) error {
	product.ID = primitive.NewObjectID()
	product.FechaActualizacion = time.Now()

	_, err := r.collection.InsertOne(ctx, &product)
	return err
}

func (r *ProductRepository) FindAll(ctx context.Context) ([]models.Product, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []models.Product
	if err = cursor.All(ctx, &products); err != nil {
		return nil, err
	}

	return products, nil
}

func (r *ProductRepository) FindByID(ctx context.Context, id string) (*models.Product, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var product models.Product
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&product)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *ProductRepository) Update(ctx context.Context, id string, product models.Product) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	product.FechaActualizacion = time.Now()

	_, err = r.collection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": product},
	)
	return err
}

func (r *ProductRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}

func (r *ProductRepository) Search(ctx context.Context, query string) ([]models.Product, error) {
	filter := bson.M{
		"$or": []bson.M{
			{"name": bson.M{"$regex": query, "$options": "i"}},
			{"licitacion": bson.M{"$regex": query, "$options": "i"}},
			{"id_convenio": bson.M{"$regex": query, "$options": "i"}},
			{"nombre_proveedor": bson.M{"$regex": query, "$options": "i"}},
			{"rut_proveedor": bson.M{"$regex": query, "$options": "i"}},
			{"id_producto": bson.M{"$regex": query, "$options": "i"}},
			{"region": bson.M{"$regex": query, "$options": "i"}},
			{"marca": bson.M{"$regex": query, "$options": "i"}},
			{"modelo": bson.M{"$regex": query, "$options": "i"}},
			{"convenio_marco": bson.M{"$regex": query, "$options": "i"}},
		},
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []models.Product
	if err = cursor.All(ctx, &products); err != nil {
		return nil, err
	}

	return products, nil
}

func (r *ProductRepository) FindAllPaginated(page, pageSize int, query bson.M) ([]*models.Product, int64, error) {
	var products []*models.Product
	ctx := context.Background()
	// Obtener el total de registros que coinciden con el filtro
	totalRecords, err := r.collection.CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	// Configurar opciones de paginación
	findOptions := options.Find()
	if page > 0 {
		findOptions.SetSkip(int64((page - 1) * pageSize))
		findOptions.SetLimit(int64(pageSize))
	}

	// Ordenar por fecha de actualización descendente
	findOptions.SetSort(bson.D{{Key: "fecha_actualizacion", Value: -1}})

	// Ejecutar la consulta
	cursor, err := r.collection.Find(ctx, query, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &products); err != nil {
		return nil, 0, err
	}

	return products, totalRecords, nil
}

func (r *ProductRepository) BuildSearchQuery(params map[string]string) bson.M {
	query := bson.M{}

	// Búsqueda por nombre
	if search := params["search"]; search != "" {
		query["$or"] = []bson.M{
			{"name": bson.M{"$regex": search, "$options": "i"}},
			{"licitacion": bson.M{"$regex": search, "$options": "i"}},
			{"id_convenio": bson.M{"$regex": search, "$options": "i"}},
			{"nombre_proveedor": bson.M{"$regex": search, "$options": "i"}},
			{"rut_proveedor": bson.M{"$regex": search, "$options": "i"}},
			{"id_producto": bson.M{"$regex": search, "$options": "i"}},
		}
	}

	// Búsqueda específica por licitación
	if licitacion := params["licitacion"]; licitacion != "" {
		query["licitacion"] = bson.M{"$regex": licitacion, "$options": "i"}
	}

	// Búsqueda por región
	if region := params["region"]; region != "" {
		query["region"] = region
	}

	// Búsqueda por marca
	if marca := params["marca"]; marca != "" {
		query["marca"] = bson.M{"$regex": marca, "$options": "i"}
	}

	// Búsqueda por convenio marco
	if convenio := params["convenio_marco"]; convenio != "" {
		query["convenio_marco"] = convenio
	}

	return query
}
