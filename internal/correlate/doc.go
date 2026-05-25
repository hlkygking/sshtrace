// Package correlate provides session correlation for sshtrace.
//
// It groups a collection of SSH sessions into clusters based on a chosen
// Strategy:
//
//   - ByUser      — all sessions belonging to the same username form one cluster.
//   - ByIP        — all sessions originating from the same source IP are grouped.
//   - ByTimeWindow — sessions whose active time ranges overlap within a
//     configurable tolerance window are placed in the same cluster.
//
// Example usage:
//
//	c, err := correlate.New(correlate.Options{
//		Strategy: correlate.ByUser,
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//	clusters := c.Group(sessions)
//	for _, cluster := range clusters {
//		fmt.Println("cluster size:", len(cluster))
//	}
package correlate
