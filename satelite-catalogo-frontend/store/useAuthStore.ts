// import { useNuxtApp, useRuntimeConfig } from "nuxt/app";
import { useRuntimeConfig } from "nuxt/app";
import { defineStore } from "pinia";

export interface AuthData {
	token: string
	expire: string
	user: {
		username: string
		_id: string
		email: string
		rut: string
		role: string[]
		created_at: string
		cc: string[]      
	}
}

async function logIn(userForm: { user: string; password: string }) {
	const { $fetchModule } = useNuxtApp()
	try {
		const dataFetch = await $fetchModule.fetchData<
			AuthData & DefaultResponse
		>({
			method: 'post',
			URL: '/auth/login',
			body: userForm,
			nuxt: false,
		})
		return dataFetch
	} catch {
		throw new Error('Credenciales inválidas')
	}
}

// async function logIn(userForm: { user: string; password: string }) {
// 	try {
// 		const runtimeConfig = useRuntimeConfig()
// 		// console.log('runtimeConfig', runtimeConfig.public.baseURL)
// 		const authData = await $fetch<AuthData>(`${runtimeConfig.public.baseURL}/auth/login`, {
// 			method: 'POST',
// 			body: userForm,
// 		})

// 		return authData
// 	} catch {
// 		throw new Error('Credenciales inválidas')
// 	}
// }

// async function logOut() {
// 	const { remove } = await useSession()
// 	await remove()
// }
// useSession es un composable que se encarga de manejar la session

// useAuthStore es un store de Pinia que se encarga de manejar la autenticación del usuario (una especie de objeto global que se puede usar en toda la app)
// se encarga de almacenar el token, la fecha de expiración y el usuario
const useAuthStore = defineStore('auth', {
	// estructura del store (datos que almacena el store)
	state: () => ({
		token: null as AuthData | null,
		expire: null as AuthData | null,
		isAuth: false,
		user: null as AuthData | null,
	}),
	// getters para obtener el estado del store
	getters: {
		getIsAuth(state) {
			return state.isAuth
		},
		getUser(state): AuthData | null {
			return state.user
		},
		getToken(state): string | null {
			return state.user?.token ?? null
		},
		getRole(state): string[] | null {
			return state.user?.user.role ?? null
		},
		getUsername(state): string | null {
			return state.user?.user.username ?? null
		},
		getID(state): string | null {
			return state.user?.user._id ?? null
		},
		getRut(state): string | null {
			return state.user?.user.rut ?? null
		},
		getEmail(state): string | null {
			return state.user?.user.email ?? null
		},
		getCC(state): string[] | null {
			return state.user?.user.cc ?? []       
		},
	},
	// acciones o metodos que se pueden ejecutar en el store
	// se pueden llamar desde cualquier parte de la app
	// y se pueden usar para modificar el estado del store
	actions: {
		async unsetAuth() {
			this.isAuth = false
			this.user = null
			this.token = null
			this.expire = null

			const { remove } = await useSession()
			await remove()
		},
		async logIn(userForm: { user: string; password: string }) {
			const dataFetch = await logIn(userForm)
			// se llama al metodo setAuth, que permite agregarlo al session de nuxt
			await this.setAuth({
				user: {
					_id: dataFetch.user._id,
					created_at: dataFetch.user.created_at,
					username: dataFetch.user.username,
					role: dataFetch.user.role,
					email: dataFetch.user.email,
					rut: dataFetch.user.rut,
					cc: dataFetch.user.cc,  
				},
				expire: dataFetch.expire,
				token: dataFetch.token,
			})
			return dataFetch
		},
		async logOut() {
			// await logOut()
			this.unsetAuth()
		},
		// aqui se setea al usuario en el store
		// se le pasa el token, la fecha de expiracion y el usuario
		// el token y la fecha de expiracion se guardan en el store
		// el usuario se guarda en el store y en la session
		async setAuth(user: AuthData) {
			this.isAuth = true
			this.user = user
			// aqui se guarda la informacion del usuario en el nuxt session
			// ya que este permite guardar la informacion del usuario en el servidor
			// y de esa forma no se pierde la sesion al recargar la pagina
			// o cambiar la pagina como ocurre en el store de pinia
			const { overwrite,session } = await useSession()
			await overwrite(user)
		},
		userRoleIs(...userTypes: string[]) {
			const roles = this.getRole
			if (!roles) return false
			return roles.some(r => userTypes.includes(r))
		},
		userRoleNotIs(...userTypes: string[]) {
			const roles = this.getRole
			if (!roles) return true
			return !roles.some(r => userTypes.includes(r))
		},
		async isUserAuthenticated() {
			const { session } = await useSession()
			if (session.value?.token){
				this.isAuth = true 
			}
			else {
				this.isAuth = false
			}
			return this.isAuth
			},
	},
})

export default useAuthStore
