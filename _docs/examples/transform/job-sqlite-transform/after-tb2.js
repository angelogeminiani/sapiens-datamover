(function () {
    const datasets = require("datasets");

    console.log("Calling after-tb2.js, LENGTH=", $data.length, JSON.stringify($variables));

    $variables.count = $variables.count || 0;
    $variables.count++;

    console.log("VARIABLES: ", JSON.stringify($variables));
    console.log("DATA 2: ", JSON.stringify($data))

    datasets.put("tb2", $data);

    // return optionally changed data
    return {
        "data": $data,
        "variables": $variables
    }
})();