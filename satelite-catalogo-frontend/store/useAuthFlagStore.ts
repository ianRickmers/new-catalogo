import { defineStore } from 'pinia'

const useAuthFlagStore = defineStore('authFlag', {
	state: () => ({
    flag: true,
	flagRole: true,
	}),
	getters: {
		
		getFlag(state) {
			return state.flag
		},
		getFlagRole(state) {
			return state.flagRole
		},
	},
	actions: {

    setFlag(authFlag: boolean){
      this.flag = authFlag
    },
	setFlagRole(authFlag: boolean){
	  this.flagRole = authFlag
	},
		
	},
})

export default useAuthFlagStore
