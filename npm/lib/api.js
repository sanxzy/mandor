/**
 * @fileoverview Mandor CLI programmatic API
 * @description Node.js API for interacting with Mandor from JavaScript/TypeScript code
 * @version 0.0.1
 */

const { spawn } = require('child_process');

/**
 * @typedef {Object} MandorOptions
 * @property {string} [cwd] - Working directory (default: process.cwd())
 * @property {boolean} [json] - Use JSON output format (default: false)
 */

/**
 * @typedef {Object} ProjectCreateOptions
 * @property {string} [name] - Project display name
 * @property {string} [goal] - Project goal description
 */

/**
 * @typedef {Object} TaskListOptions
 * @property {string} [project] - Filter by project ID
 * @property {string} [feature] - Filter by feature ID
 * @property {string} [status] - Filter by status (pending, ready, in_progress, done, blocked, cancelled)
 * @property {string} [priority] - Filter by priority (P0-P5)
 */

/**
 * Mandor CLI wrapper class for programmatic access
 * @class
 * @version 0.0.1
 * @example
 * const mandor = new Mandor({ cwd: '/path/to/project', json: true });
 * const tasks = await mandor.taskList({ project: 'api', status: 'pending' });
 */
class Mandor {
  /**
   * Creates a new Mandor instance
   * @constructor
   * @param {MandorOptions} [options] - Configuration options
   * @example
   * const mandor = new Mandor({ cwd: '/my/project', json: true });
   */
  constructor(options = {}) {
    /** @type {string} Working directory */
    this.cwd = options.cwd || process.cwd();
    /** @type {boolean} Use JSON output */
    this.json = options.json || false;
  }

  /**
   * Initializes a new Mandor workspace
   * @async
   * @param {string} name - Workspace name
   * @returns {Promise<number>} Exit code from the CLI
   * @throws {Error} If initialization fails
   * @example
   * await mandor.init('My AI Project');
   */
  async init(name) {
    return this._run('init', [name]);
  }

  /**
   * Creates a new project
   * @async
   * @param {string} id - Project identifier (lowercase, hyphens only)
   * @param {ProjectCreateOptions} [options] - Project options
   * @returns {Promise<number>} Exit code from the CLI
   * @example
   * await mandor.projectCreate('api', {
   *   name: 'API Service',
   *   goal: 'Implement REST API with user management'
   * });
   */
  async projectCreate(id, options = {}) {
    const args = ['project', 'create', id];
    if (options.name) args.push('--name', options.name);
    if (options.goal) args.push('--goal', options.goal);
    return this._run(...args);
  }

  /**
   * Lists tasks with optional filters
   * @async
   * @param {TaskListOptions} [options] - Filter options
   * @returns {Promise<Object[]|number>} Task list (JSON) or exit code
   * @example
   * const tasks = await mandor.taskList({
   *   project: 'api',
   *   status: 'pending',
   *   json: true
   * });
   */
  async taskList(options = {}) {
    const args = ['task', 'list'];
    if (options.project) args.push('--project', options.project);
    if (options.feature) args.push('--feature', options.feature);
    if (options.status) args.push('--status', options.status);
    if (options.priority) args.push('--priority', options.priority);
    if (options.json) args.push('--json');
    return this._run(...args);
  }

  /**
   * Gets detailed information about a task
   * @async
   * @param {string} taskId - Task identifier
   * @returns {Promise<Object|number>} Task details (JSON) or exit code
   * @example
   * const task = await mandor.taskDetail('api-feature-abc-task-xyz789', { json: true });
   */
  async taskDetail(taskId) {
    return this._run('task', 'detail', taskId);
  }

  /**
   * Updates a task's status or metadata
   * @async
   * @param {string} taskId - Task identifier
   * @param {Object} updates - Fields to update
   * @returns {Promise<number>} Exit code from the CLI
   * @example
   * await mandor.taskUpdate('api-feature-abc-task-xyz789', {
   *   status: 'in_progress'
   * });
   */
  async taskUpdate(taskId, updates = {}) {
    const args = ['task', 'update', taskId];
    if (updates.status) args.push('--status', updates.status);
    if (updates.name) args.push('--name', updates.name);
    if (updates.priority) args.push('--priority', updates.priority);
    return this._run(...args);
  }

  /**
   * Lists features with optional filters
   * @async
   * @param {Object} [options] - Filter options
   * @returns {Promise<Object[]|number>} Feature list (JSON) or exit code
   * @example
   * const features = await mandor.featureList({ project: 'api', json: true });
   */
  async featureList(options = {}) {
    const args = ['feature', 'list'];
    if (options.project) args.push('--project', options.project);
    if (options.status) args.push('--status', options.status);
    if (options.json) args.push('--json');
    return this._run(...args);
  }

  /**
   * Lists issues with optional filters
   * @async
   * @param {Object} [options] - Filter options
   * @returns {Promise<Object[]|number>} Issue list (JSON) or exit code
   * @example
   * const issues = await mandor.issueList({ project: 'api', type: 'bug', json: true });
   */
  async issueList(options = {}) {
    const args = ['issue', 'list'];
    if (options.project) args.push('--project', options.project);
    if (options.type) args.push('--type', options.type);
    if (options.status) args.push('--status', options.status);
    if (options.json) args.push('--json');
    return this._run(...args);
  }

  /**
   * Gets workspace status and statistics
   * @async
   * @param {Object} [options] - Options
   * @returns {Promise<Object|number>} Status (JSON) or exit code
   * @example
   * const status = await mandor.status({ json: true });
   */
  async status(options = {}) {
    const args = ['status'];
    if (options.json) args.push('--json');
    return this._run(...args);
  }

  /**
   * Internal method to run mandor CLI commands
   * @private
   * @async
   * @param {...string} args - Command arguments
   * @returns {Promise<Object|number>} Result based on json option
   * @throws {Error} If process fails
   */
  _run(...args) {
    return new Promise((resolve, reject) => {
      const proc = spawn('mandor', args, {
        cwd: this.cwd,
        stdio: this.json ? 'pipe' : 'inherit'
      });

      if (this.json) {
        let data = '';
        proc.stdout.on('data', chunk => data += chunk);
        proc.on('close', (code) => {
          try {
            resolve(JSON.parse(data));
          } catch {
            resolve(data);
          }
        });
      } else {
        proc.on('close', resolve);
      }
    });
  }
}

module.exports = Mandor;
