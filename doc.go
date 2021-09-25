/*
package skipfilter provides a data structure that combines a skiplist with a roaring bitmap cache.

This library was created to efficiently filter a multi-topic message input stream against a set of subscribers,
each having a list of topic subscriptions expressed as regular expressions. Idealy, each subscriber should test
each topic at most once to determine whether it wants to receive messages from the topic.

In this case, the skip list provides an efficient discontinuous slice of subscribers and the roaring bitmap for each
topic provides an efficient ordered discontinuous set of all subscribers that have indicated that they wish to
receive messages on the topic.

Filter bitmaps are stored in an LRU cache of variable size (default 100,000).

This package is not theadsafe.
*/
package skipfilter
