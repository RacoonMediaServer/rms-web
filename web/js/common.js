function copyToClipboard(me,id) {
    var elem = document.getElementById(id);

    window.getSelection().selectAllChildren(elem);
    document.execCommand("Copy")

    me.value = "Скопировано"
}