var modal = document.getElementById('myModal');
var openModalButton = document.getElementById('openModalButton');
var closeModalSpan = document.getElementsByClassName('close')[0];
var confirmButton = document.getElementById('confirmButton');
var cancelButton = document.getElementById('cancelButton');
// Close the modal when the close span is clicked
closeModalSpan.onclick = function() {
    modal.style.display = 'none';
}
// Close the modal when the cancel button is clicked
cancelButton.onclick = function() {
    modal.style.display = 'none';
}
// Confirm action when the confirm button is clicked
confirmButton.onclick = function() {
    modal.style.display = 'none';
    uploadFile()
}

var del_all_id = document.getElementById('del_all_id');

var dropZone = document.getElementById('code_div');//code_div drop_zone
var fileList = document.getElementById('file_list');
var fileInput = document.getElementById('fileInput');
var upfiles = new Array();
dropZone.addEventListener('dragover', function(e) {
    e.preventDefault();
    dropZone.classList.add('dragover');
});

dropZone.addEventListener('dragleave', function(e) {
    e.preventDefault();
    dropZone.classList.remove('dragover');
});

dropZone.addEventListener('drop', function(e) {
    e.preventDefault();
    dropZone.classList.remove('dragover');
    var files = e.dataTransfer.files;
    handleFiles(files);
});

function handleFiles(files) {
    for (var i = 0; i < files.length; i++) {
        var file = files[i];
        var listItem = document.createElement('li');
        listItem.textContent = file.name + ' - ' + formatBytes(file.size);
        fileList.appendChild(listItem);
        upfiles.push(file)
    }
    modal.style.display = 'block';
}

function formatBytes(bytes, decimals = 2) {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const dm = decimals < 0 ? 0 : decimals;
    const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' ' + sizes[i];
}

var authcode = ''
function init() {
    authcode = localStorage.getItem('password');
    if (authcode){
        auth(authcode)
    }else{
        document.getElementById('content').style.display = 'none';
        document.getElementById('auth').style.display = 'block';
    }
    var input = document.getElementById('passwordInput')
    input.textContent = authcode
    input.addEventListener('keyup', function(event) {
        if (event.key === 'Enter') {
            var password = event.target.value;
            auth(password)
            event.target.value = ''; // Clear the input field
        }
    });

    GetConfig((data) => {
        if (data.AppVersion !== ''){
            console.log('====>data.AppVersion ',data.AppVersion)
            document.title = "Go文件上传助手v" + data.AppVersion;
            title = document.getElementById('title_h2')
            title.textContent = 'Go文件上传助手v' + data.AppVersion + ' ';
            var aname = document.createElement("a");
            aname.target = '_blank'
            aname.id = 'pubip'
            title.appendChild(aname)
            getBuildInfo(data)
        }
    })
}

function getBuildInfo(jsonData) {
    const DOUBLE_CLICK_TIME = 300;
    let lastClickTime = 0;
    document.getElementById('title_h2').addEventListener('click', function(event) {
        // 获取当前时间
        const now = new Date().getTime();
        // 检查当前点击与上次点击的时间差
        if (now - lastClickTime < DOUBLE_CLICK_TIME) {
            // 如果两次点击的时间差小于阈值，则触发双击事件
            //handleDoubleClick(event);
            build_info = '应用名称：' + jsonData.AppName
            build_info += '\r\n应用版本：' + jsonData.AppVersion
            build_info += '\r\n编译版本：' + jsonData.BuildVersion
            build_info += '\r\n编译日期：' + jsonData.BuildTime
            build_info += '\r\nGitRevision：' + jsonData.GitRevision
            build_info += '\r\nGitBranch：' + jsonData.GitBranch
            build_info += '\r\nGoVersion：' + jsonData.GoVersion
            Toast(build_info,10)
        } else {
            // 否则，更新上次点击的时间
            lastClickTime = now;
        }
    });
}
function cache(){
    localStorage.setItem('key', 'value');
    var value = localStorage.getItem('key');
    localStorage.removeItem('key');
}
function copyToClipboard(text) {
    // Create a textarea element
    var textarea = document.createElement("textarea");
    textarea.value = text;
    // Append the textarea to the document body
    document.body.appendChild(textarea);
    // Select the textarea's content
    textarea.select();
    // Execute the copy command
    document.execCommand("copy");
    // Remove the textarea from the document body
    document.body.removeChild(textarea);
    showToast('已复制：' + text)
}
function copy1(event) {
    event.stopPropagation();
    copyToClipboard(document.getElementById("up_1").textContent)
}
function copy2(event) {
    event.stopPropagation();
    copyToClipboard(document.getElementById("up_2").textContent)
}
function copy3(event) {
    event.stopPropagation();
    copyToClipboard(document.getElementById("up_3").textContent)
}
function createcode(token) {
    document.getElementById('code_div').style.display = 'block';

    const code3 = document.getElementById('up_3')
    code3.textContent = `bash <(curl -s -S -L ${window.location.origin}/up) `

    const code2 = document.getElementById('up_2')
    code2.textContent = `curl -H "Authorization: ${token}" -F "file=@/root/a.log" -F "file=@/root/b.log" ${window.location.origin}/upload`

    const code1 = document.getElementById('up_1')
    code1.textContent = `curl -H "Authorization: ${token}" -F "file=@/root/a.log" ${window.location.origin}/upload`
}

function showToast(content) {
    Toast(content,3)
}

function Toast(content,timeout) {
    var toastElement = document.getElementById("toast");
    // 设置Toast文本
    toastElement.innerText = content;
    // 显示Toast
    toastElement.style.display = "block";
    // 3秒后隐藏Toast
    setTimeout(function () {
        toastElement.style.display = "none";
        document.getElementById('progress').style.width = '0%';
        document.getElementById('progress').textContent = '';
        document.getElementById('progressBar').style.display = "none";
    }, 1000*timeout);
}

function showFiles(path) {
    authcode = localStorage.getItem('password');
    var xhr = new XMLHttpRequest();
    var url = '/upload';
    url += `?origin=${window.location.origin}`
    if (path){
        url += `?origin=${window.location.origin}&path=${path}`
    }
    xhr.open('GET', url, true);
    xhr.setRequestHeader("Authorization",authcode)
    console.log('url',url);
    console.log('authcode',authcode);
    xhr.onreadystatechange = function () {
        if (xhr.readyState === 1) {
            // 在这里处理loading状态，例如显示loading动画
            console.log('Loading...');
            showLoading('正在获取文件清单，请稍等～')
        }else if (isHttpOk(xhr)) {//xhr.readyState === 4 &&
            filejson = JSON.parse(xhr.response)
            if (filejson.code === 0){
                console.log('showFiles',xhr.response)
                var table = document.getElementById("myTable");
                var tbody = table.getElementsByTagName("tbody")[0];
                tbody.innerHTML = '';
                if (filejson.data){
                    // 使用 for...of 循环倒序遍历数组
                    for (var element of filejson.data.reverse()) {
                        addItemByGet(element)
                    }
                }
            }else{
                console.log('失败了',filejson.msg)
            }
            hideLoading()
        } else {
            // 请求失败或还未完成
            //console.error('get files err ',xhr.response);
            console.log('get files err ',xhr.readyState,xhr.status,xhr.response);
        }
    };

    xhr.send();
}

function getPubIp() {
    var xhr = new XMLHttpRequest();
    var url = '/getip';
    xhr.open('GET', url, true);
    xhr.onreadystatechange = function () {
        if (xhr.readyState === 1) {
            // 在这里处理loading状态，例如显示loading动画
            console.log('Loading...');
            showLoading('正在获取文件清单，请稍等～')
        }else if (isHttpOk(xhr)) {//xhr.readyState === 4 &&
            //console.log('getPubIp...',xhr);
            // 获取<a>标签的引用
            var link = document.getElementById('pubip');
            // 设置超链接的目标URL
            link.href = xhr.responseText;
            // 设置链接文本
            link.textContent = '公网';
            hideLoading()
        } else {
            // 请求失败或还未完成
            console.log('getPubIp err ',xhr.readyState,xhr.status);
        }
    };

    xhr.send();
}

function clearTable() {
    // 获取表格对象
    var table = document.getElementById("myTable");
    // 获取表格主体
    var tbody = table.getElementsByTagName("tbody")[0];
    // 移除表格主体中的所有行
    while (tbody.firstChild) {
        tbody.removeChild(tbody.firstChild);
    }
    //getFiles()
}

function auth(password) {
    var xhr = new XMLHttpRequest();
    xhr.open('POST', '/auth', true);
    xhr.setRequestHeader("Authorization",password)
    xhr.onreadystatechange = function() {
        console.log('====',xhr.readyState,xhr.status)
        if (xhr.status === 200){
            document.getElementById('content').style.display = 'block';
            document.getElementById('auth').style.display = 'none';
            localStorage.setItem('password', password);
            showFiles(null)
            getPubIp()
            console.log('sucess',xhr.status,xhr.responseText)
            showToast('认证成功')
            createcode(password)
        }else{
            console.log('failed',xhr.status)
            showToast('认证失败')
            document.getElementById('content').style.display = 'none';
            document.getElementById('auth').style.display = 'block';
        }
    };
    xhr.send();
}

function GetConfig(callback) {
    var xhr = new XMLHttpRequest();
    xhr.open('POST', '/config', true);
    //xhr.setRequestHeader("Authorization",password)
    xhr.onreadystatechange = function() {
        console.log('====',xhr.readyState,xhr.status)
        if (isHttpOk(xhr)){
            console.log('sucess',xhr.status,xhr.responseText)
            filejson = JSON.parse(xhr.response)
            if (filejson.code === 0){
                callback(filejson.data)
            }
        }
    };
    xhr.send();
}

function uploadFiles(formData,total_size){
    var xhr = new XMLHttpRequest();

    const progressBar = document.getElementById('progressBar');
    const speedElement = document.getElementById('speed');
    const formatSpeed = (bytesPerSecond) => {
        const kiloBytesPerSecond = bytesPerSecond / (1024*1024);
        return kiloBytesPerSecond.toFixed(2) + ' MB/s';
    };// 记录上传开始时间

    let startTime;
    let startBytes = 0;
    startTime = new Date().getTime();

    // 监听进度事件
    xhr.upload.addEventListener('progress', function (event) {
        if (event.lengthComputable) {
            var percentComplete = (event.loaded / event.total) * 100;
            document.getElementById('progress').style.width = percentComplete + '%';
            var roundedResult = percentComplete.toFixed(1);
            document.getElementById('progress').textContent = roundedResult + '%';
        }

        const currentTime = new Date().getTime();
        const elapsedSeconds = (currentTime - startTime) / 1000;
        const uploadedBytes = event.loaded;
        const speed = (uploadedBytes - startBytes) / elapsedSeconds;

        // 更新上传速度
        const speedText = formatSpeed(speed);
        speedElement.textContent = `上传速度: ${speedText}`;

        // 更新进度条
        const progress = (uploadedBytes / total_size) * 100;
        progressBar.value = progress;

        // 更新起始时间和字节数
        startTime = currentTime;
        startBytes = uploadedBytes;

    });

    xhr.onreadystatechange = function () {
        if (xhr.readyState === 1) {
            // 在这里处理loading状态，例如显示loading动画
            console.log('Loading...');
            showLoading('正在上传文件～')
        }else if (xhr.readyState === 4) {
            if (isHttpOk(xhr)) {
                // 文件上传成功
                console.log('File uploaded successfully!');
                console.log(xhr)
                filejson = JSON.parse(xhr.response)
                if (filejson.code === 0){
                    console.log('成功了哦')
                    if (filejson.data){
                        for (var item of  filejson.data){
                            addItemByUpload(item)
                        }
                    }
                    showToast('上传成功')
                    // document.getElementById('fileInput').textContent = '';
                    //localStorage.setItem("token", token);
                    //console.log(token)
                    document.getElementById('progress').textContent = '上传成功';
                    while (fileList.firstChild) {
                        fileList.removeChild(fileList.firstChild);
                    }
                }else{
                    console.log('失败了',filejson.msg)
                    showToast('调用失败:' + filejson.msg)
                }
            } else {
                // 文件上传失败
                console.error('File upload failed',xhr);
                showToast('文件上传失败，请重新上传',xhr.status,xhr.statusText)
            }

            // 关闭模态框
            // closeUploadModal();
            clearFileInput()
            hideLoading()
            while (fileList.firstChild) {
                fileList.removeChild(fileList.firstChild);
            }
        }
    };
    progressBar.style.display = "block"
    // 将 '/upload' 替换为服务器端处理文件上传的路径
    xhr.open('POST', '/upload', true);
    xhr.setRequestHeader("Authorization",authcode)
    xhr.setRequestHeader("source","web")
    xhr.send(formData);
}

function uploadFile(){
    var filecount = upfiles.length
    if (filecount == 0){
        showToast('请选择文件上传～')
    }else{
        var formData = new FormData();
        var total_size = 0;
        for (var i = 0; i < filecount; i++) {
            var file = upfiles.pop()
            formData.append('file', file);
            total_size += file.size
        }
        console.log('total_size',total_size)
        formData.append('token', authcode);
        uploadFiles(formData,total_size)
    }
}


function deletefile(files,callback) {
    const jsonData = {
        files: files,
    };
    const xhr = new XMLHttpRequest();
    xhr.open('DELETE', '/upload', true);
    xhr.setRequestHeader("Authorization",authcode)
    xhr.setRequestHeader('Content-Type', 'application/json');
    xhr.onreadystatechange = function () {
        if (xhr.readyState === 1){
            console.log('Loading...');
            showLoading('正在删除文件～')
        } else if (xhr.readyState === 4) {
            if (isHttpOk(xhr)) {
                console.log('Post deleted successfully', xhr.responseText);
                filejson = JSON.parse(xhr.response)
                if (filejson.code === 0){
                    data = filejson.data
                    if (data.length == 1){
                        if (data[0].sucess){
                            console.log('成功了哦')
                            callback(true,filejson.msg)
                        }
                    }else{
                        callback(true,data)
                    }
                }else{
                    callback(false,filejson.msg)
                }

            } else {
                console.error('Request failed');
                callback(false,'' + xhr.status)
            }
            hideLoading()
        }
    };
    xhr.send(JSON.stringify(jsonData));
}

function clearFileInput() {
    // 获取文件输入字段
    var fileInput = document.getElementById('fileInput');

    // 创建一个新的文件输入字段
    var newFileInput = document.createElement('input');
    newFileInput.type = 'file';
    newFileInput.id = 'fileInput';

    // 将新的文件输入字段替换原有的文件输入字段
    fileInput.parentNode.replaceChild(newFileInput, fileInput);
}

function formatFileSize(sizeInBytes) {
    const kilobyte = 1024;
    const megabyte = kilobyte * 1024;
    const gigabyte = megabyte * 1024;

    if (sizeInBytes < kilobyte) {
        return sizeInBytes.toFixed(2) + ' B';
    } else if (sizeInBytes < megabyte) {
        return (sizeInBytes / kilobyte).toFixed(2) + ' KB';
    } else if (sizeInBytes < gigabyte) {
        return (sizeInBytes / megabyte).toFixed(2) + ' MB';
    } else {
        return (sizeInBytes / gigabyte).toFixed(2) + ' GB';
    }
}

function insertRow(tbody,newRow,newItem) {
    //<td><input type="checkbox" class="selectRow"></td>
    var cell0 = newRow.insertCell(0);
    var cell1 = newRow.insertCell(1);
    // var cell2 = newRow.insertCell(2);
    var cell2 = newRow.insertCell(2);
    var cell3 = newRow.insertCell(3);
    var cell4 = newRow.insertCell(4);
    var input = document.createElement('input');
    input.type = 'checkbox'
    input.className = 'selectRow'
    input.alt = newItem.path
    // 添加 'change' 事件的事件监听器
    input.addEventListener('change', function () {
        // 检查复选框是否被选中
        if (this.checked) {
            del_all_id.style.display = 'block';
        }
    });


    //cell1.innerHTML = "<a href="+newItem.path+">"+ newItem.name +"</a>";
    var aname = document.createElement("a");
    aname.textContent = newItem.name;
    aname.href = newItem.path;
    aname.target = '_blank'

    var copylinkbtn = document.createElement('button');
    copylinkbtn.textContent = '复制';
    copylinkbtn.style = 'margin-right: 5px;'
    copylinkbtn.addEventListener('click', function () {
        let text = window.origin + newItem.path
        copyToClipboard(text)
        showToast('已复制：' + text)
    });

    var downloadbtn = document.createElement('button');
    downloadbtn.style = 'margin-right: 5px; margin-left: 5px;'
    downloadbtn.textContent = '下载';
    downloadbtn.addEventListener('click', function () {
        //aname.click()
        var url = encodeURIComponent(newItem.path);
        //window.open(url, '_blank');
        var link = document.createElement("a");
        link.textContent = newItem.name;
        link.href = url;
        link.download = newItem.name;
        link.click()
    });
    // var encodedPath = encodeURIComponent(newItem.path);
    // cell2.innerHTML = "<a href="+encodedPath+">下载</a>";

    // 创建按钮并设置事件处理程序
    var delbtn = document.createElement('button');
    delbtn.textContent = '删除';
    delbtn.style = 'margin-right: 5px;'
    //delbtn.className = 'delete-btn'
    delbtn.addEventListener('click', function () {
        // 当按钮点击时触发的事件
        var result = window.confirm(newItem.name + " 确定要删除这个文件吗？");
        if (result) {
            deletefile([newItem.path], function (ok, msg) {
                if (ok) {
                    showToast('删除成功' + newItem.path)
                    tbody.removeChild(newRow);
                } else {
                    showToast('删除失败 ' + msg)
                }
            })
        } else {
        }
    });




    cell0.appendChild(input);
    cell1.appendChild(aname);
    cell2.appendChild(downloadbtn);
    cell2.appendChild(delbtn);
    cell2.appendChild(copylinkbtn);
    cell3.innerHTML = formatFileSize(newItem.size)
    cell4.innerHTML = newItem.modTime;
    // cell4.appendChild(operate);
    //console.log('==>', newItem)
}

function addItemByUpload(newItem) {
    var table = document.getElementById("myTable");
    var tbody = table.getElementsByTagName("tbody")[0];
    var newRow = tbody.insertRow(0);
    insertRow(tbody,newRow,newItem)
}

function addItemByGet(newItem) {
    var table = document.getElementById("myTable");
    var tbody = table.getElementsByTagName("tbody")[0];
    var newRow = tbody.insertRow();
    insertRow(tbody,newRow,newItem)
}

function selectAllRows() {
    var checkboxes = document.getElementsByClassName('selectRow');
    var selectAllCheckbox = document.getElementById('selectAll');

    for (var i = 0; i < checkboxes.length; i++) {
        checkboxes[i].checked = selectAllCheckbox.checked;
    }
    showToast('您的总文件数量为：'+checkboxes.length)
    // 检查复选框是否被选中
    if (selectAllCheckbox.checked) {
        del_all_id.style.display = 'block';
    } else {
        del_all_id.style.display = 'none';
    }
}


function deleteSelectedRows() {
    var result = window.confirm("确定要删除这个文件吗？");
    if (result) {
        var checkboxes = document.getElementsByClassName('selectRow');
        var dataArr = [];
        for (var i = checkboxes.length - 1; i >= 0; i--) {
            if (checkboxes[i].checked) {
                var row = checkboxes[i].parentNode.parentNode;
                dataArr.push(checkboxes[i].alt)
                row.parentNode.removeChild(row);
            }
        }
        deletefile(dataArr,function (ok,msg) {
            showToast('成功删除文件数：'+msg.length)
        })
    } else {
        showToast('感谢大爷不删之恩～')
    }

}

function showLoading(msg) {
    // 显示loading状态
    // document.getElementById('overlay').style.display = 'block';
    document.getElementById('overlay').style.display = 'flex';
    document.getElementById('loading').innerText = msg
}

function hideLoading() {
    // 隐藏loading状态
    document.getElementById('overlay').style.display = 'none';
    document.getElementById('loading').innerText = '加载中...'
}

function isHttpOk(xhr) {
    if (xhr.status === 200 && xhr.response && xhr.response.length > 0){
        return true;
    }
    return false
}

init()