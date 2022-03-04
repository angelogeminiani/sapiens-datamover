(function () {
    console.log("CONTEXT");
    console.log("$data", JSON.stringify($data));

    return {
        "data": $data,
        "variables": $variables
    }
})();