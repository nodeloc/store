/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      fontFamily: {
        sans: ['Inter', 'system-ui', '-apple-system', 'sans-serif'],
        mono: ['JetBrains Mono', 'monospace']
      },
      colors: {
        zinc: {
          50: '#fafafa',
          100: '#f4f4f5',
          200: '#e4e4e7',
          300: '#d4d4d8',
          400: '#a1a1aa',
          500: '#71717a',
          600: '#52525b',
          700: '#3f3f46',
          800: '#27272a',
          900: '#18181b',
          950: '#09090b'
        },
        brand: {
          green: '#009966',
          orange: '#FF9933',
          'green-light': '#00b377',
          'orange-light': '#ffad5c',
        }
      },
      backgroundImage: {
        'brand-gradient': 'linear-gradient(135deg, #009966 0%, #FF9933 100%)',
        'brand-gradient-r': 'linear-gradient(135deg, #FF9933 0%, #009966 100%)',
        'brand-gradient-subtle': 'linear-gradient(135deg, rgba(0,153,102,0.08) 0%, rgba(255,153,51,0.08) 100%)',
        'brand-gradient-hover': 'linear-gradient(135deg, rgba(0,153,102,0.12) 0%, rgba(255,153,51,0.12) 100%)',
      },
      boxShadow: {
        'card': '0 1px 3px rgba(0,0,0,0.04), 0 1px 2px rgba(0,0,0,0.06)',
        'card-hover': '0 10px 25px rgba(0,0,0,0.08), 0 4px 10px rgba(0,0,0,0.04)',
        'glass': '0 4px 30px rgba(0,0,0,0.05)',
        'glow': '0 0 20px rgba(0,153,102,0.15)',
      },
      animation: {
        'fade-in': 'fadeIn 0.5s ease-out',
        'slide-up': 'slideUp 0.4s ease-out',
        'float': 'float 6s ease-in-out infinite',
      },
      keyframes: {
        float: {
          '0%, 100%': { transform: 'translateY(0)' },
          '50%': { transform: 'translateY(-10px)' },
        }
      }
    },
  },
  plugins: [],
}
