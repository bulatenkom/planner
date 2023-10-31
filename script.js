"use strict";

document.querySelector("#add-new-event button").addEventListener("click", (e) => {

    let data = `
        {
            "title": "${document.querySelector("#add-new-event input[name=title]").value}",
            "description": "${document.querySelector("#add-new-event input[name=description]").value}",
            "duration": "${document.querySelector("#add-new-event input[name=duration]").value}",
            "plannedOn": "${document.querySelector("#add-new-event input[name=plannedOn]").value}",
            "type": "${document.querySelector("#add-new-event select[name=type]").value}",
            "done": ${document.querySelector("#add-new-event input[name=done]").checked}
        }
    `
    // 'done' param is not used on EVENT creation, but can be useful in PUT

    console.log(data)

    fetch("/events", {
        method: "POST",
        body: data
    })
})