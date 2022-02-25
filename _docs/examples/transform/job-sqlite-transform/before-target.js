(function () {
        try {
            const datasets = require("datasets");

            console.log("Calling before-target.js, LENGTH=", JSON.stringify($data), JSON.stringify($variables));

            $variables.count = $variables.count || 0;
            $variables.count++;

            console.log("VARIABLES: ", JSON.stringify($variables));

            // alter data
            const tb1 = datasets.get("tb1");
            const tb2 = datasets.get("tb2");
            if (!!tb1 && !!tb2) {
                let data = [];
                console.log("TABLES", JSON.stringify(tb1), JSON.stringify(tb2));
                if (tb1.length === tb2.length) {
                    for (let i = 0; i < tb1.length; i++) {
                        const d1 = tb1[i];
                        const d2 = tb2[i];
                        const item = {id:d1["id"], name: d1["name"], surname: d2["surname"]};
                        data.push(item);
                    }
                    console.log("TARGET", JSON.stringify(data));
                    datasets.put("target", data); // just for debug
                    return {
                        "data": data, // returns new data for next SQL command
                        "variables": $variables
                    }
                }
            } else {
                console.log("TABLES ARE NULL");
            }
        } catch
            (err) {
            console.error("before-target.js error: ", err);
        }

        // return optionally changed data
        return {
            "data": $data,
            "variables": $variables
        }
    }
)();