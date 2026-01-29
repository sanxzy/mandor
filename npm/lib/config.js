/**
 * @fileoverview Mandor configuration management
 * @description Reads and writes Mandor configuration using .mandorrc.json
 * @version 0.0.1
 */

const path = require('path');
const fs = require('fs');

/** @type {string} Default config filename */
const CONFIG_FILENAME = '.mandorrc.json';

/**
 * Configuration management class for Mandor projects
 * @class
 * @version 0.0.1
 * @example
 * const config = new MandorConfig('/path/to/project');
 * const defaultPriority = config.get('priority.default', 'P3');
 * config.set('theme', 'dark');
 */
class MandorConfig {
  /**
   * Creates a new MandorConfig instance
   * @constructor
   * @param {string} projectRoot - Root directory of the project
   * @example
   * const config = new MandorConfig('/my/project');
   */
  constructor(projectRoot) {
    /** @type {string} Path to config file */
    this.configPath = path.join(projectRoot, CONFIG_FILENAME);
    /** @type {string} Project root directory */
    this.projectRoot = projectRoot;
  }

  /**
   * Gets a configuration value
   * @param {string} key - Configuration key (supports dot notation, e.g., 'priority.default')
   * @param {*} [defaultValue] - Default value if key not found
   * @returns {*} The configuration value or default
   * @example
   * const priority = config.get('priority.default', 'P3');
   * const strict = config.get('strictMode', false);
   */
  get(key, defaultValue = undefined) {
    if (!fs.existsSync(this.configPath)) return defaultValue;
    const config = JSON.parse(fs.readFileSync(this.configPath, 'utf-8'));

    // Support dot notation for nested keys
    const keys = key.split('.');
    let value = config;
    for (const k of keys) {
      if (value && typeof value === 'object' && k in value) {
        value = value[k];
      } else {
        return defaultValue;
      }
    }
    return value !== undefined ? value : defaultValue;
  }

  /**
   * Sets a configuration value
   * @param {string} key - Configuration key (supports dot notation)
   * @param {*} value - Value to set
   * @returns {void}
   * @example
   * config.set('priority.default', 'P2');
   * config.set('theme', 'dark');
   */
  set(key, value) {
    const config = fs.existsSync(this.configPath)
      ? JSON.parse(fs.readFileSync(this.configPath, 'utf-8'))
      : {};

    // Support dot notation for nested keys
   .split('.');
    const keys = key let current = config;
    for (let i = 0; i < keys.length - 1; i++) {
      const k = keys[i];
      if (!current[k] || typeof current[k] !== 'object') {
        current[k] = {};
      }
      current = current[k];
    }
    current[keys[keys.length - 1]] = value;

    fs.writeFileSync(this.configPath, JSON.stringify(config, null, 2));
  }

  /**
   * Deletes a configuration key
   * @param {string} key - Configuration key to delete
   * @returns {boolean} True if key was deleted
   * @example
   * config.delete('theme');
   */
  delete(key) {
    if (!fs.existsSync(this.configPath)) return false;
    const config = JSON.parse(fs.readFileSync(this.configPath, 'utf-8'));

    const keys = key.split('.');
    let current = config;
    for (let i = 0; i < keys.length - 1; i++) {
      if (!current[keys[i]]) return false;
      current = current[keys[i]];
    }

    if (current[keys[keys.length - 1]] !== undefined) {
      delete current[keys[keys.length - 1]];
      fs.writeFileSync(this.configPath, JSON.stringify(config, null, 2));
      return true;
    }
    return false;
  }

  /**
   * Checks if a configuration key exists
   * @param {string} key - Configuration key to check
   * @returns {boolean} True if key exists
   * @example
   * if (config.has('theme')) { console.log('Theme is set'); }
   */
  has(key) {
    return this.get(key) !== undefined;
  }

  /**
   * Gets all configuration as an object
   * @returns {Object} Full configuration object
   * @example
   * const allConfig = config.getAll();
   */
  getAll() {
    if (!fs.existsSync(this.configPath)) return {};
    return JSON.parse(fs.readFileSync(this.configPath, 'utf-8'));
  }

  /**
   * Clears all configuration
   * @returns {void}
   * @example
   * config.clear();
   */
  clear() {
    if (fs.existsSync(this.configPath)) {
      fs.unlinkSync(this.configPath);
    }
  }
}

module.exports = MandorConfig;
