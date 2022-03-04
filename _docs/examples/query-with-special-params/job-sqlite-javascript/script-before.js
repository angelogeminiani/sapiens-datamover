(function () {
    console.log("BEFORE");


    $variables["array"] = ["A", "B"];
    console.log("$variables", JSON.stringify($variables));

    return {
        "data": $data,
        "variables": $variables
    }
})();