<?php

class Validator {
export const AuthRules = {
	password: {
		message: {
			en: {
				number: '%s must include one or multiple numbers',
				lower: '%s must include one or multiple lowercase characters',
				upper: '%s must include one or multiple uppercase characters',
				special: '%s must include one or multiple special characters',
			},
		},
		validate(input, withErrors = false) {
			const errors = [];

			if (!/\d/.test(input)) {
				errors.push('regex.number');
			}

			if (!/[a-z]/.test(input)) {
				errors.push('regex.lower');
			}

			if (!/[A-Z]/.test(input)) {
				errors.push('regex.upper');
			}

			if (!/\W/.test(input)) {
				errors.push('regex.special');
			}

			if (input.length < 8) {
				errors.push('min');
			}

			if (input.length > 130) {
				errors.push('max');
			}

			if (withErrors) {
				return {
					errors: errors,
					valid: errors.length === 0,
				}
			}

			return errors.length === 0
		},
	},
}