function UploadFile(){
    let input = document.querySelector('#file-upload input[type=file]');
    if (input.files.length > 0){
        const fileName = document.querySelector('#file-upload .file-name');
        fileName .textContent = input.files[0].name;

        // set content for zone 
        const display = document.getElementById("input-doc");
        console.log(input.files[0].name)
        console.log("here we are");
        
        const reader = new FileReader();
        reader.onload = function fileReadCompleted() {
          // when the reader is done, the content is in reader.result.
          console.log(reader.result);
          display.value =reader.result
        };
        reader.readAsText(input.files[0]);
    }
}

function SendToBackend(){
    console.log("there you are")
    // validating type file 
    const fileName = document.querySelector('#file-upload .file-name')
    const extension = fileName.textContent.match(/\.[0-9a-z]+$/i);
    const display = document.getElementById("input-doc");
    const uri = "/api/consulize";

    // looks awful
    if (extension){ 
        (extension[0] === '.json' || extension[0] === '.hcl') ? Post(extension[0], uri, display.value): window.alert("not supported extension");
        return
    }


    window.alert("Supported file extensions are json/hcl");
}

function Post(extension, uri, payload){
    var xhr = new XMLHttpRequest();
    xhr.open("POST",uri, true);
    xhr.setRequestHeader("Content-Type", "application/json");

    xhr.onreadystatechange = function(){
        if (xhr.readyState === 4){
            console.log("......... sent ZZZ");
            console.log(xhr.responseText);
        }
    }

    var data = JSON.stringify({
            "extension": extension,
            "payload" : window.btoa(payload)
        }
    );

    xhr.send(data);
}