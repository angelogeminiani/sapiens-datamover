(function () {
    const datasets = require("datasets");

    console.log("Calling after-tb1.js, LENGTH=", $data.length, JSON.stringify($variables));

    $variables.count = $variables.count || 0;
    $variables.count++;

    console.log( "VARIABLES: ", JSON.stringify($variables) )
    console.log( "DATA 1: ", JSON.stringify($data) )

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

    // return optionally changed data
    return {
        "data": $data,
        "variables": $variables
    }
})();