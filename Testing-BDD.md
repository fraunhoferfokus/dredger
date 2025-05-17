# BDD Generation

The input feature file has to have the following structure in order for the generator to create the godog test file

-   ```Scenario: Test GET Request for url <"regex of the url">
      When I send GET request to <"actual endpoint"> with payload <"payload that needs to be sent">
      Then The response for url <"endpoint again"> with request method <"request method"> should be <status code>
    ```
-   So an example for a scenario would be
-   ```
      Scenario: Test GET Request for url "/store/{id}"
      When I send GET request to "/store/100" with payload ""
      Then The response for url "/store/100" with request method "GET" should be 404
    ```
-   After you have created a feature file that follows this structure, you can generate the godog file by running `dredger generate-bdd <path to the file>`
