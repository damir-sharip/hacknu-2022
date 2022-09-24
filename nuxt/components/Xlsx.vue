<template>
    <div>
        <b-form-file
            @change="handleFile"
            placeholder="Choose file"
            drop-placeholder="Drop file here..."
        ></b-form-file>
    </div>
</template>

<script>
import * as XLSX from "xlsx";

export default {
    methods: {
        handleFile(e) {
            const files = e.target.files;
            if (!files.length) return;
            const file = files[0];
            const reader = new FileReader();
            reader.onload = (e) => {
                const data = e.target.result;
                const workbook = XLSX.read(data, { type: "binary" });

                const resultingData = {};

                workbook.SheetNames.forEach((sheetName) => {
                    const worksheet = workbook.Sheets[sheetName];
                    const data = XLSX.utils.sheet_to_json(worksheet);

                    const dataSnakeCase = data.map((row) => {
                        const newRow = {};
                        Object.keys(row).forEach((key) => {
                            const newKey = key
                                .split(" ")
                                .map((word) => word.toLowerCase())
                                .join("_");
                            newRow[newKey] = row[key];
                        });
                        return newRow;
                    });

                    const dataNull = dataSnakeCase.map((row) => {
                        const newRow = {};
                        Object.keys(row).forEach((key) => {
                            newRow[key] = row[key] === "null" ? null : row[key];
                            if (key === "identifier" && newRow[key] === null) {
                                newRow[key] = "UNKNOWN";
                            }
                        });
                        return newRow;
                    });

                    const dataFiltered = dataNull.filter((row) => {
                        return row.latitude && row.longitude;
                    });
                    resultingData[sheetName] = dataFiltered;
                });

                console.log(resultingData);

                this.$store.commit("routes/SET_ROUTES", resultingData);
            };
            reader.readAsBinaryString(file);
        },
    },
};
</script>