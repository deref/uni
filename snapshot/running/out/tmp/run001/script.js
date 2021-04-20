require('source-map-support').install();
const { main } = require('./bundle.js');
if (typeof main === 'function') {
	const args = process.argv.slice(2);
	void (async () => {
		const exitCode = await main(...args);
		process.exit(exitCode ?? 0);
	})();
} else {
	process.stderr.write('error: /current/working/path/snapshot/running/exit.ts does not export a main function\n', () => {
		process.exit(1);
	});
}
