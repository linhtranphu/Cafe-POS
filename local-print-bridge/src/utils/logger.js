const LOG_LEVEL = process.env.LOG_LEVEL || 'info'

const levels = {
  debug: 0,
  info: 1,
  warn: 2,
  error: 3
}

const currentLevel = levels[LOG_LEVEL] || levels.info

function formatMessage(level, ...args) {
  const timestamp = new Date().toISOString()
  const levelStr = level.toUpperCase().padEnd(5)
  return `[${timestamp}] [${levelStr}] ${args.join(' ')}`
}

function debug(...args) {
  if (currentLevel <= levels.debug) {
    console.log(formatMessage('debug', ...args))
  }
}

function info(...args) {
  if (currentLevel <= levels.info) {
    console.log(formatMessage('info', ...args))
  }
}

function warn(...args) {
  if (currentLevel <= levels.warn) {
    console.warn(formatMessage('warn', ...args))
  }
}

function error(...args) {
  if (currentLevel <= levels.error) {
    console.error(formatMessage('error', ...args))
  }
}

module.exports = {
  debug,
  info,
  warn,
  error
}
