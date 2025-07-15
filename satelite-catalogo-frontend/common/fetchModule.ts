/* eslint-disable no-console */
/* eslint-disable security-node/detect-crlf */
// Modules
import { v4 as uuidv4 } from 'uuid'
// Types
import type { FetchError } from 'ofetch'
import { capitalizeFirstLetter } from '../utils/format'
import { useRoute, useRouter, useRuntimeConfig } from 'nuxt/app'

enum HTTPMethods {
	'post',
	'put',
	'get',
	'delete',
}

export interface ConfigFetch {
	method: keyof typeof HTTPMethods
	URL: string

	body?: BodyInit | null | undefined | Record<string, any>
	spinnerStatus?: boolean
	headers?: HeadersInit
	token?: string | null
	nuxt?: boolean
	// Abort is for all same URLs
	abort?: {
		url?: string
		onChangePath?: boolean
	}
	params?: { [x: string]: unknown }
	signal?: AbortSignal
	responseType?: 'blob' | 'text' | 'json' | 'stream' | 'arrayBuffer'
	scopeSpinner?: string
}

export type DefaultResponse = {
	success: boolean
	message?: string
}

export type ErrorFetch = {
	success: boolean
	message: string
	statusCode: number
}

export const ERROR_ABORT = 'The user aborted a request'

export class Fetch {
	private publicApi: string | undefined
	// string: URL, Array: currentFetchController
	private currentFetch: Map<
		string,
		Array<{
			id: string
			controller: AbortController
			path: string
			config: ConfigFetch
		}>
	> = new Map()

	private counters: {
		counterFetch: number
		counterGetFetch: number
	} = new Proxy(
		{ counterFetch: 0, counterGetFetch: 0 },
		{
			set(obj, prop, value) {
				if (prop === 'counterFetch' || prop === 'counterGetFetch') {
					obj[prop] = value

					if (prop === 'counterFetch' && value === 0)
						useSpinner().value = false
					if (prop === 'counterGetFetch' && value === 0)
						useSpinnerGet().value = false
				}
				return true
			},
		},
	)

	private readonly spinner = useSpinner()
	private readonly spinnerGet = useSpinnerGet()
	private readonly scopeSpinner = useScopeSpinner()
	// Route
	private path = useRoute().path

	constructor() {
		this.watchPath()
	}

	private generateFetchId(): string {
		return `fetch_id_${uuidv4()}`
	}

	private isFetchError(error: unknown): error is FetchError {
		return typeof error === 'object' && error !== null && 'message' in error
	}

	private watchPath() {
		useRouter().beforeEach((to) => {
			this.path = to.path

			this.currentFetch.forEach((fetchController, url) => {
				fetchController.forEach(({ id, controller, path }) => {
					if (
						path !== this.path &&
						controller.signal &&
						!controller.signal.aborted
					) {
						const config = this.currentFetch
							.get(url)
							?.find((f) => f.id === id)
						if (config?.config.abort?.onChangePath)
							controller.abort()
					}
				})
			})

			return true
		})
	}

	/**
	 * Handles errors in fetch request
	 * @param error The error to handle
	 * @param save Indicates wheter the error should be saved in the error store
	 * @returns An ErrorFetch object with information about the error
	 */
	handleError(error: unknown): ErrorFetch {
		let errorFetch: ErrorFetch
		if (this.isFetchError(error)) {
			const message = error.data?.message ?? error.message
			errorFetch = {
				success: false,
				message: capitalizeFirstLetter(message),
				statusCode: error.response?.status ?? 500,
			}
		} else if (error instanceof Error) {
			errorFetch = {
				success: false,
				message: capitalizeFirstLetter(error.message),
				statusCode: 500,
			}
		} else {
			errorFetch = {
				success: false,
				message: 'fetch',
				statusCode: 500,
			}
		}
		return errorFetch
	}

	private getFetch({
		method,
		body,
		nuxt,
		token,
		responseType,
		signal,
		headers,
	}: ConfigFetch) {
		if (!this.publicApi) {
			const config = useRuntimeConfig()
			this.publicApi = (config.public.backBaseUrl as string) ?? ''
		}
		return $fetch.create({
			baseURL: !nuxt ? this.publicApi : '',
			parseResponse: responseType === 'blob' ? undefined : JSON.parse,
			responseType,
			headers: {
				...(token ? { Authorization: `Bearer ${token}` } : {}),
				...headers,
			},

			method,
			body,
			signal,
			onRequestError({ request, error }) {
				// Log error
				console.error(`[fetch request error] ${request} ${error}`)
			},
			onResponseError({ request, response }) {
				// Log response
				console.error('[fetch response error]', request, response.body)
			},
			mode: 'cors',
		})
	}

	/**
	 * Removes a fetch from the list of current requests
	 * @param id The ID of the fetch to remove
	 * @param counters Counters to reduce in one for GET and non-GET requests
	 * @param key The key of the fetch request in the currentFetch map. (URL)
	 * @param scopeSpinner Optional. The scope of the spinner
	 */
	private popFetch(
		id: string,
		counters: { get: boolean; fetch: boolean },
		key: string,
		scopeSpinner?: string,
	) {
		if (counters.get) this.counters.counterGetFetch -= 1
		if (counters.fetch) this.counters.counterFetch -= 1

		const index = this.currentFetch.get(key)?.findIndex((f) => f.id === id)
		this.currentFetch.get(key)?.splice(index ?? 0, 1)
		if (scopeSpinner) this.scopeSpinner.value.delete(scopeSpinner)
	}

	/**
	 * Performs a fetch request and returns the response
	 * @param config The configuration of the fetch request
	 * @returns A promise with the response from the fetch request
	 */
	async fetchData<T = any>(config: ConfigFetch): Promise<T> {
		// Add Params
		if (config.params) {
			let hasQuery = config.URL.includes('?')
			for (const [key, value] of Object.entries(config.params).filter(
				([_, value]) => value !== undefined,
			)) {
				config.URL += `${hasQuery ? '&' : '?'}${key}=${value}`
				hasQuery = true
			}
		}

		const key = config.URL.split('?')[0]
		// Abort all fetchs
		const abortKey = config.abort?.url === 'same' ? key : config?.abort?.url
		if (
			config.abort &&
			config.abort?.url &&
			this.currentFetch.has(abortKey ?? '')
		) {
			const struct = this.currentFetch.get(abortKey ?? '')
			struct?.forEach((c) => {
				c.controller.abort()
			})
			struct?.splice(0, struct.length)
		}
		// Id
		const id = this.generateFetchId()
		// Controller signal
		const controller = new AbortController()
		config.signal = controller.signal

		if (this.currentFetch.has(key)) {
			this.currentFetch.get(key)?.push({
				id,
				controller,
				path: this.path,
				config,
			})
		} else {
			this.currentFetch.set(key, [
				{ id, controller, path: this.path, config },
			])
		}
		// Create fetch
		const apiFetch = this.getFetch(config)

		// Scope fetch
		if (!config.scopeSpinner) config.scopeSpinner = 'default'
		this.scopeSpinner.value.set(config.scopeSpinner, true)

		// Methods
		if (config.method !== 'get' || config.spinnerStatus) {
			this.spinner.value = true
			this.counters.counterFetch += 1
		}
		if (config.method === 'get') {
			this.spinnerGet.value = true
			this.counters.counterGetFetch += 1
		}

		const dataFetch = await apiFetch<T & DefaultResponse>(config.URL).catch(
			(err) => {
				this.popFetch(
					id,
					{
						get: config.method === 'get',
						fetch:
							config.method !== 'get' ||
							(config.spinnerStatus ?? false),
					},
					key,
					config.scopeSpinner,
				)
				throw err
			},
		)
		this.popFetch(
			id,
			{
				get: config.method === 'get',
				fetch:
					config.method !== 'get' || (config.spinnerStatus ?? false),
			},
			key,
			config.scopeSpinner,
		)
		return dataFetch as T & DefaultResponse
	}
}