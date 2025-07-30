package constants

import "time"

const PER_PERSIST_BATCH_SIZE = 10

const PERSIST_PENDING_CHECK_INTERVAL = 5 * time.Second

const MIN_NUM_OF_WORKERS = 1

const HIGH_PENDING_COUNT = 50
const LOW_PENDING_COUNT = 20

const RESTORE_LIMIT = 10

const CLIENT_IDLE_TIMEOUT_CHECK_INTERVAL = 5 * time.Second
const CLIENT_IDLE_TIMEOUT = 3 * time.Minute
