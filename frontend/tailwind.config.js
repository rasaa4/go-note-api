/** @type {import('tailwindcss').Config} */
export default {
  content: ["./index.html", "./src/**/*.{js,ts,jsx,tsx}"],
  theme: {
    extend: {
      colors: {
        paper: {
          DEFAULT: "#F4F6F5",
          card: "#FFFFFF",
          line: "#DCE3E0",
        },
        ink: {
          DEFAULT: "#12201C",
          soft: "#3F544D",
          faint: "#78897F",
        },
        moss: {
          50: "#EEF5F1",
          100: "#D7E8DE",
          300: "#8FBBA3",
          500: "#2F6E52",
          600: "#255A43",
          700: "#1C4634",
        },
        rust: {
          500: "#C15B3A",
          600: "#A6492C",
        },
        amber: {
          500: "#C98A2C",
        },
      },
      fontFamily: {
        display: ["'Newsreader'", "serif"],
        sans: ["'IBM Plex Sans'", "system-ui", "sans-serif"],
        mono: ["'IBM Plex Mono'", "monospace"],
      },
      boxShadow: {
        card: "0 1px 2px rgba(18,32,28,0.06), 0 1px 0 rgba(18,32,28,0.04)",
      },
    },
  },
  plugins: [],
};
