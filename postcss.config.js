import tailwind from "@tailwindcss/postcss";
import autoprefixer from "autoprefixer";
import cssnanoPlugin from "cssnano";
export default {
    plugins: [
        autoprefixer(),
        tailwind(),
        // cssnanoPlugin()
    ]
}