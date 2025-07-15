import type { ConfigFetch, Fetch } from '../common/fetchModule'
import useAuthStore from '../store/useAuthStore'


export abstract class Service {
	private readonly fetchModule: Fetch
	protected readonly authStore = useAuthStore()
	// private readonly toastsStore = useToast()

	constructor(fetch: Fetch) {
		this.fetchModule = fetch
	}

	protected fetch<T = any>(
		config: Omit<ConfigFetch, 'token'>,
		omitToken = false,
	) {
		return this.fetchModule.fetchData<T>({
			...config,
			token: omitToken ? undefined : this.authStore.getToken,
		})
	}

	private handleError(error: unknown) {
		return this.fetchModule.handleError(error)
	}

}