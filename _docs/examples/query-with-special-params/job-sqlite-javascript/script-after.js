(function () {
    console.log("AFTER");
    console.log("$data", JSON.stringify($data));

    return {
        "data": $data,
        "variables": $variables
    }
})();