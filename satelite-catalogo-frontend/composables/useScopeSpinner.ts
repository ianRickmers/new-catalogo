import { useState } from "nuxt/app";

export default function () {
	return useState('scopespinner', () => new Map<string, boolean>())
}