package node

type Node interface {
	/**
	 * Get incoming request per minute ({@code pass + block}).
	 *
	 * @return total request count per minute
	 */
	RequestInMinute() uint64

	/**
	 * Get pass count per minute.
	 *
	 * @return total passed request count per minute
	 */
	PassInMinute() uint64

	/**
	 * Get {@link Entry#exit()} count per minute.
	 *
	 * @return total completed request count per minute
	 */
	SuccessInMinute() uint64

	/**
	 * Get blocked request count per minute (totalBlockRequest).
	 *
	 * @return total blocked request count per minute
	 */
	BlockInMinute() uint64

	/**
	 * Get exception count per minute.
	 *
	 * @return total business exception count per minute
	 */
	ErrorInMinute() uint64

	/**
	 * Get pass request per second.
	 *
	 * @return QPS of passed requests
	 */
	PassQps() float64

	/**
	 * Get block request per second.
	 *
	 * @return QPS of blocked requests
	 */
	BlockQps() float64

	/**
	 * Get {@link #passQps()} + {@link #blockQps()} request per second.
	 *
	 * @return QPS of passed and blocked requests
	 */
	TotalQps() float64

	/**
	 * Get {@link Entry#exit()} request per second.
	 *
	 * @return QPS of completed requests
	 */
	SuccessQps() float64

	/**
	 * Get estimated max success QPS till now.
	 *
	 * @return max completed QPS
	 */
	MaxSuccessQps() float64

	/**
	 * Get exception count per second.
	 *
	 * @return QPS of exception occurs
	 */
	ErrorQps() float64

	/**
	 * Get average rt per second.
	 *
	 * @return average response time per second
	 */
	AvgRtInSecond() float64

	/**
	 * Get average rt per second.
	 *
	 * @return average response time per second
	 */
	AvgRtInMinute() float64

	/**
	 * Get minimal response time in second
	 *
	 * @return recorded minimal response time
	 */
	MinRtInSecond() uint64

	/**
	 * Get minimal response time in minute
	 *
	 * @return recorded minimal response time
	 */
	MinRtInMinute() uint64

	/**
	 * Get current active goroutine count.
	 *
	 * @return current active goroutine count
	 */
	CurGoroutineNum() uint64

	/**
	 * Get last second block QPS.
	 */
	PreviousBlockQps() uint64

	/**
	 * Last window QPS.
	 */
	PreviousPassQps() uint64

	/**
	 * Fetch all valid metric nodes of resources.
	 *
	 * @return valid metric nodes of resources
	 */
	Metrics() map[uint64]*MetricNode

	/**
	 * Add pass count.
	 *
	 * @param count count to add pass
	 */
	AddPassRequest(count uint64)

	/**
	 * Add rt and success count.
	 *
	 * @param rt      response time
	 * @param success success count to add
	 */
	AddRtAndSuccess(rt uint64, success uint64)

	/**
	 * Increase the block count.
	 *
	 * @param count count to add
	 */
	AddBlockRequest(count uint64)

	/**
	 * Add the biz error count.
	 *
	 * @param count count to add
	 */
	AddErrorRequest(count uint64)

	/**
	 * Increase current thread count.
	 */
	IncreaseGoroutineNum()

	/**
	 * Decrease current thread count.
	 */
	DecreaseGoroutineNum()

	/**
	 * Reset the internal counter.
	 */
	Reset()
}
