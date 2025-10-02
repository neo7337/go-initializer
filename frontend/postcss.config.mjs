import postcssCalc from 'postcss-calc';

export default {
  plugins: {
    "@tailwindcss/postcss": {},
    "postcss-calc": postcssCalc(),
  }
}