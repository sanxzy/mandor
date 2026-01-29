/**
 * @fileoverview Mandor CLI npm package entry point
 * @description Main export for programmatic usage and CLI access
 * @version 0.0.1
 */

const Mandor = require('./api');
const MandorConfig = require('./config');
const { resolve, listCachedBinaries, clearCache } = require('./resolve');
const { downloadBinary, getCurrentPlatform } = require('./download');
const { install } = require('./install');

/**
 * Main export object containing all public APIs
 * @namespace mandor
 * @version 0.0.1
 * @example
 * const mandor = require('@mandor/cli');
 *
 * // CLI access via npx or npm script
 * // $ npx mandor init "My Project"
 *
 * // Programmatic usage
 * const cli = new mandor.Mandor({ json: true });
 * await cli.init('My Project');
 */
module.exports = {
  /**
   * Mandor CLI wrapper class for programmatic access
   * @type {typeof Mandor}
   * @memberof mandor
   * @example
   * const mandor = require('@mandor/cli');
   * const cli = new mandor.Mandor({ json: true });
   */
  Mandor,

  /**
   * Configuration management class
   * @type {typeof MandorConfig}
   * @memberof mandor
   * @example
   * const config = new mandor.MandorConfig('/path/to/project');
   * const priority = config.get('priority.default', 'P3');
   */
  MandorConfig,

  /**
   * Resolves the Mandor binary path
   * @function
   * @param {Object} [options] - Resolution options
   * @param {string} [options.version] - Version to use
   * @returns {Promise<string>} Path to binary
   * @memberof mandor
   * @example
   * const binaryPath = await mandor.resolve({ version: 'latest' });
   */
  resolve,

  /**
   * Lists all cached binary versions
   * @function
   * @returns {Object[]} Cached binary info
   * @memberof mandor
   * @example
   * const cached = mandor.listCachedBinaries();
   */
  listCachedBinaries,

  /**
   * Clears all cached binaries
   * @function
   * @returns {number} Number of binaries removed
   * @memberof mandor
   * @example
   * const removed = mandor.clearCache();
   */
  clearCache,

  /**
   * Downloads a Mandor binary
   * @function
   * @param {string} version - Version to download
   * @param {string} [platform] - Target platform
   * @param {string} [arch] - Target architecture
   * @returns {Promise<string>} Path to downloaded binary
   * @memberof mandor
   * @example
   * const binary = await mandor.downloadBinary('1.0.0', 'darwin', 'x64');
   */
  downloadBinary,

  /**
   * Gets current platform information
   * @function
   * @returns {{platform: string, arch: string}} Platform info
   * @memberof mandor
   * @example
   * const { platform, arch } = mandor.getCurrentPlatform();
   */
  getCurrentPlatform,

  /**
   * Runs post-install setup
   * @function
   * @param {Object} [options] - Install options
   * @returns {Promise<string>} Path to installed binary
   * @memberof mandor
   * @example
   * await mandor.install({ version: 'latest' });
   */
  install,

  // Version info
  /** @type {string} Package version */
  version: '0.0.1',

  /** @type {string} Supported Mandor version range */
  mandorVersionRange: '>=0.0.1'
};

// CLI entry point when bin/mandor is executed
if (require.main === module) {
  resolve()
    .then(binaryPath => {
      const { spawn } = require('child_process');
      const args = process.argv.slice(2);
      const proc = spawn(binaryPath, args, {
        stdio: 'inherit',
        cwd: process.cwd()
      });
      proc.on('exit', process.exit);
    })
    .catch(error => {
      console.error('Failed to start Mandor:', error.message);
      process.exit(1);
    });
}
