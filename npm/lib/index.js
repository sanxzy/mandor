/**
 * @fileoverview Mandor CLI npm package entry point
 * @description Main export for CLI access
 * @version 0.0.2
 */

const { install } = require('./install');
const MandorConfig = require('./config');

module.exports = {
  MandorConfig,
  install,
  version: '0.0.2'
};

if (require.main === module || process.env.npm_lifecycle_event === 'postinstall') {
  install().catch(error => {
    console.error('Failed to install Mandor:', error.message);
    process.exit(1);
  });
}
