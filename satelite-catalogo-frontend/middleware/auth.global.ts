import { navigateTo, useRequestHeaders } from 'nuxt/app'
import { UserTypes } from '~/models/user/user.model'

type UserTypesKeys = keyof typeof UserTypes

export default defineNuxtRouteMiddleware(async (to) => {
  if (process.server) return
// MANEJO DE AUTENTICACION
  const authStore = useAuthStore()
  let isAuth = authStore.isAuth

  if (!isAuth) {
    try {
      const session = await $fetch('/api/session', {
        headers: useRequestHeaders(['cookie'])
      })

      if (session?.token) {
        await authStore.setAuth({
          token: session.token,
          expire: session.expire,
          user: {
            _id: session.user._id,
            username: session.user.username,
            created_at: session.user.created_at,
            email: session.user.email,
            rut: session.user.rut,
            cc: session.user.cc,
            role: session.user.role, // <--- array de roles
          },
        })
        isAuth = true
      }
    } catch (error) {
      console.error('Error fetching session:', error)
      isAuth = false
    }
  }

  const authRequired = to.meta?.auth ?? true
  if (authRequired && !isAuth) {
    return window.location.href = '/login'
  }
// MANEJO DE ROLES
  const currentPath = to.path
  const roles = (to.meta.userType ?? []) as UserTypesKeys[]
  if (roles.length > 0) {
    const userRoles: string[] = authStore.getRole || []
    console.log('Roles de la ruta:', roles)
    console.log('Roles del usuario:', userRoles)
    console.log('Ruta actual:', currentPath)
    const requiredRoles = roles.map((r) => UserTypes[r])

    const hasRole = userRoles.some((userRole) => requiredRoles.includes(userRole))
    console.log('Tiene rol requerido:', hasRole)
    if (!hasRole && currentPath !== '/login') {
      // Elegimos la home del usuario según sus roles

      let homePage = "/login"

      if (userRoles.includes('Usuario') || userRoles.includes('Administrador')) {
        homePage = "/solicitudes"
      } else {
        await authStore.logOut()
        homePage = "/login"
      }

      if (currentPath !== homePage) {
        return window.location.href = homePage
      }
    }
  }
})
