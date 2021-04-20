require('source-map-support').install();

const { inspect } = require('util');
process.on('uncaughtException', (exception) => {
  process.stderr.write('uncaught exception: ' + inspect(exception) + '\n', () => {
    process.exit(1);
  });
});
process.on('unhandledRejection', (reason, promise) => {
  process.stderr.write(
    'unhandled rejection at: ' + inspect(promise) + '\nreason: ' + inspect(reason) + '\n',
    () => {
      process.exit(1);
    },
  );
})

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
