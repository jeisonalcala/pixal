import svelte from 'rollup-plugin-svelte';
import resolve, { nodeResolve } from '@rollup/plugin-node-resolve';
import commonjs from '@rollup/plugin-commonjs';
import livereload from 'rollup-plugin-livereload';
import terser from '@rollup/plugin-terser';
import babel from '@rollup/plugin-babel'

const production = !process.env.ROLLUP_WATCH;

const builddir = "../res/static/svelte"
export default [
	"file_viewer",
	"filesystem",
	"user_home",
	"user_file_manager",
	"admin_panel",
	"home_page",
	"text_upload",
	"speedtest",
].map((name, index) => ({
	input: `src/${name}.js`,
	output: {
		sourcemap: true,
		format: 'iife',
		name: 'app',
		file: `${builddir}/${name}.js`,
	},
	plugins: [
		svelte({
			compilerOptions: {
				// enable run-time checks when not in production
				dev: !production,
				// we'll extract any component CSS out into
				// a separate file - better for performance
				// css: css => {
				// 	css.write(`${name}.css`);
				// },
			},
			emitCss: false,
		}),

		babel({
			extensions: [".js", ".html", ".svelte"],
			babelHelpers: "bundled",
		}),

		// If you have external dependencies installed from
		// npm, you'll most likely need these plugins. In
		// some cases you'll need additional configuration -
		// consult the documentation for details:
		// https://github.com/rollup/plugins/tree/master/packages/commonjs
		resolve({
			browser: true,
			exportConditions: ['svelte'],
			extensions: ['.svelte'],
		}),
		commonjs(),
		nodeResolve(),

		// In dev mode, call `npm run start` once
		// the bundle has been generated
		// !production && serve(),

		// Watch the `public` directory and refresh the
		// browser on changes when not in production
		!production && livereload({
			watch: `${builddir}/${name}.*`,
			port: 5000 + index,
		}),

		// If we're building for production (npm run build
		// instead of npm run dev), minify
		production && terser(),
	],
	watch: {
		clearScreen: false
	},
}));
