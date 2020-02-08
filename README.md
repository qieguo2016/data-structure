# 数据结构(Go实现)

## [跳跃表(非并发安全)](./skip_list/skip_list.go)

![skip list](./resource/800px-Skip_list.svg.png)

增加了向前指针的链表叫作跳表。跳表全称叫做跳跃表，简称跳表。跳表是一个随机化的数据结构，实质就是一种可以进行二分查找的有序链表。跳表在原有的有序链表上面增加了多级索引，通过索引来实现快速查找。跳表不仅能提高搜索性能，同时也可以提高插入和删除操作的性能。

这里采用redis底层类似的实现，在每层上增加了偏移量的记录，好处是在按排行取元素的时候可以先从上层按偏移量快速定位到目标位置，不需要在底层链表进行遍历定位。

Reference:
1. https://en.wikipedia.org/wiki/Skip_list


## [B+树(非并发安全)](./tree/b_tree.go)
 
![b+ tree](./resource/800px-Bplustree.png)
 
B+ 树是一种树数据结构，是一个n叉树，每个节点通常有多个孩子，一颗B+树包含根节点、内部节点和叶子节点。根节点可能是一个叶子节点，也可能是一个包含两个或两个以上孩子节点的节点。 
B+ 树通常用于数据库和操作系统的文件系统中。 NTFS, ReiserFS, NSS, XFS, JFS, ReFS 和BFS等文件系统都在使用B+树作为元数据索引。 
B+ 树元素自底向上插入，其特点是能够保持数据稳定有序，其插入与修改拥有较稳定的对数时间复杂度。 

Reference:
1. https://en.wikipedia.org/wiki/Skip_list


## [缓存(lru/lfu)](./cache/lru.go)

LRU是Least Recently Used的缩写，即最近最少使用，是一种常用的页面置换算法，选择最近最久未使用的页面予以淘汰。该算法赋予每个页面一个访问字段，用来记录一个页面自上次被访问以来所经历的时间 t，当须淘汰一个页面时，选择现有页面中其 t 值最大的，即最近最久未使用的页面予以淘汰。  
LFU（least frequently used (LFU) page-replacement algorithm）。即最不经常使用页置换算法，要求在页置换时置换引用计数最小的页，因为经常使用的页应该有一个较大的引用次数。但是有些页在开始时使用次数很多，但以后就不再使用，这类页将会长时间留在内存中，因此可以将引用计数寄存器定时右移一位，形成指数衰减的平均使用次数。

- LRU cache [LRU cache](./cache/lru.go)
- LFU cache [LFU cache](./cache/lfu.go)

Reference:
1. https://en.wikipedia.org/wiki/Cache_replacement_policies#Least_recently_used_(LRU) 
2. https://en.wikipedia.org/wiki/Least_frequently_used


## [堆](./heap/heap.go)

堆是一种带有顺序结构的完全二叉树，分为大根堆和小根堆，根据完全二叉和父子大小关系，利用数组结构比较容易实现堆结果。
golang源码中也实现了一个小根堆(代码在container/heap/)，采用接口化的设计，实用性大大提升，值得好好学习一番，主要亮点：
1. 接口化设计，只要实现heap接口即可使用
2. 复用sort接口实现，最大程度复用
3. 采用循环代推递归实现调整 

## [链表](./list/dual_link_list.go)

golang实现的单链表和双链表结构和源码分析。
golang源码的双向链表实现(代码在container/list/)亮点：
1. 双向链表为环形结构，前后指针调整方便
2. 节点元素与链表分开两种数据结构

## [并发实现](./concurrent/concurrent_example.go)

- [once实现](./concurrent/once.go)
- [无锁队列链表实现](./concurrent/linked_queue.go)
