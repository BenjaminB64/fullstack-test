import { nextui } from '@nextui-org/react'

/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
    "./node_modules/@nextui-org/theme/dist/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {},
  },
  darkMode: "class",
  plugins: [nextui(
      {
        layout: {
            fontSize: {
                tiny: "0.85rem", // text-tiny
                small: "0.975rem", // text-small
                medium: "1.1rem", // text-medium
                large: "1.225rem", // text-large
            },
        },
        themes: {
          light: {
            colors: {
                primary: '#484EF4',
            }
          },
          dark: {
            colors: {
                primary: '#484EF4',
            }
          }
        }
      }
  )],
}

