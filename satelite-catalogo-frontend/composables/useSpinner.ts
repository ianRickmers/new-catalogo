import { useState } from "nuxt/app";

export default function () {
	return useState('spinner', () => false)
}