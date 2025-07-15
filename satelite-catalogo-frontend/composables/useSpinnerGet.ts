import { useState } from "nuxt/app";

export default function () {
	return useState('spinnerGet', () => false)
}