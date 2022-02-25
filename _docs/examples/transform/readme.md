# Sample for Data Transformation #

In this example we run three actions:

 - `source-tb1`: SELECT all items from tb1 Table
 - `source-tb2`: SELECT all items from tb2 Table
 - `target`: INSERT context data into target Database

Take a look at scripts assigned to each action to understand how to use
a custom module named `datasets`

```javascript
(function () {
    const datasets = require("datasets");

    // add context data to datasets storage component
    datasets.put("tb1", $data);

    // test dataset has been added
    const tb1 = datasets.get("tb1");
    if (!tb1){
        throw "PROBLEM ADDING DATASET";
    }
    // test a map() method
    const selected = datasets.map("tb1", function(item){
        console.log("   --->  SELECTING ITEM: ", JSON.stringify(item));
        return item; // select
    });
    console.log("SELECTED with datasets.map(): ", JSON.stringify(selected));
    // test for loop
    const exitMessage = datasets.for("tb1", function(item){
        console.log("   --->  READING ITEM: ", JSON.stringify(item));
    });

   
    return {
        "data": $data,
        "variables": $variables
    }
})();
```