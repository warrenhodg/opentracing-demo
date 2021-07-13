# opentracing-demo
1. Run the jaeger server by running run-server.sh.
2. Run client1, and client2
3. curl http://localhost:8080/foo or /bar to generate spans from client1 -> client2 
4. curl http://localhost:8081/foor or /bar to generate spans from client2 -> client1
5. You should see the distributed traces in jaeger by accessing UI at http://http://localhost:16686
