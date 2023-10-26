"use strict";

document.querySelector("#add-new-slot button").addEventListener("click", (e) => {

    let data = `
        {
            "title": "${document.querySelector("#add-new-slot input[name=title]").value}",
            "description": "${document.querySelector("#add-new-slot input[name=description]").value}",
            "duration": "${document.querySelector("#add-new-slot input[name=duration]").value}",
            "plannedOn": "${document.querySelector("#add-new-slot input[name=plannedOn]").value}",
            "type": "${document.querySelector("#add-new-slot select[name=type]").value}",
            "done": ${document.querySelector("#add-new-slot input[name=done]").checked}
        }
    `
    // 'done' param is not used on SLOT creation, but can be useful in PUT

    console.log(data)

    fetch("/slots", {
        method: "POST",
        body: data
    })
})