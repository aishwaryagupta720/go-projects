# Execution Command

The way to execute my code is:  

``` bash
go run main.go
```
open load_balancing_comparison.html to see the generated simulation

# Extra Algorithms

## 1. Round Robin
- **Algorithm Name:** Round Robin
  - **Explanation:**
    - Round Robin works by sequentially assigning each incoming request to the next server in the list. When it reaches the end of the list, it starts over from the beginning. This ensures an even distribution of requests across all servers.
  - **Comparison with Power of Two in a distributed environment:**
    - Round Robin performs worse than Power of Two Choices in a distributed environment because it does not consider the current load on each server. If one server becomes slower or fails, Round Robin continues to send requests to it, potentially causing delays. In contrast, Power of Two Choices tends to balance the load more effectively by selecting the least loaded of two randomly chosen servers.

## 2. Power of Two Choices
- **Algorithm Name:** Power of Two Choices
  - **Explanation:**
    - The Power of Two Choices algorithm randomly selects two servers and assigns the incoming request to the server with the lower load. This method helps in distributing the requests more evenly and reduces the chances of overloading a single server.
  - **Comparison with Power of Two in a distributed environment:**
    - Power of Two Choices is designed to perform well in a distributed environment by balancing the load more dynamically and efficiently compared to other static methods like Round Robin.

## 3. Random Allocation
- **Algorithm Name:** Random Allocation
  - **Explanation:**
    - Random Allocation assigns each incoming request to a server chosen randomly from available servers. This method is simple and easy to implement.
  - **Comparison with Power of Two in a distributed environment:**
    - Random Allocation performs worse than Power of Two Choices because it does not consider the current load on the servers. This can lead to uneven distribution of requests and potential overloading of some servers.

## 4. Least Recently Contacted
- **Algorithm Name:** Least Recently Contacted
  - **Explanation:**
    - Least Recently Contacted works by selecting the server that has not been assigned a request for the longest time. This method ensures that all servers are utilized evenly over time.
  - **Comparison with Power of Two in a distributed environment:**
    - Least Recently Contacted performs worse than Power of Two Choices in a distributed environment because it does not consider the current load on each server. While it helps in utilizing servers evenly, it may lead to overloading slower servers.

## 5. Weighted Round Robin
- **Algorithm Name:** Weighted Round Robin
  - **Explanation:**
    - Weighted Round Robin works similarly to Round Robin but assigns more requests to servers with higher weights. This method takes into account the different capacities of servers, distributing the load according to their weights.
  - **Comparison with Power of Two in a distributed environment:**
    - Weighted Round Robin performs worse than Power of Two Choices because it does not dynamically consider the current load on the servers. Although it distributes requests based on server capacity, it may not adapt quickly to changing conditions in a distributed environment.

## 6. Weighted Random
- **Algorithm Name:** Weighted Random
  - **Explanation:**
    - Weighted Random assigns requests to servers randomly, with the probability of selecting a server proportional to its weight. This method accounts for the different capacities of servers.
  - **Comparison with Power of Two in a distributed environment:**
    - Weighted Random performs worse than Power of Two Choices because it does not dynamically balance the load. While it considers server capacity, it may lead to uneven distribution if server loads change frequently.



# Grading

The grading is based on the completion of the following criteria _and_ your ability
to explain your code. I suggest you leave many comments that explain the what and why of the code
so you're prepared for when I ask you about it. 

You are required to complete each criteria in this assignment. 

| Points | ID     | Test Criteria                                                                        |
| -----: | ------ | ------------------------------------------------------------------------------------ |
|      5 | CHARTS | Your code generates charts (a long with instructions on how to regenerate them)      |
|      5 | PO2    | Create an experiment that pits the "Power of 2 Random Choices" against "Round Robin" |

You can include other load balancing algorithms, as long as you explain the algorithm. Each additional algorithm and explanation will be 5 extra points. Although the explanation is about distributed load balancing, your implementation can assume a single load balancer.
