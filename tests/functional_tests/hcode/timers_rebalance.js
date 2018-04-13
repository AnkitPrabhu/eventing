function OnUpdate(doc, meta) {
    var expiry = Math.round((new Date()).getTime() / 1000) + 5;
    docTimer(timerCallback, meta.id, expiry);
}
function OnDelete(meta) {
    var expiry = Math.round((new Date()).getTime() / 1000) + 5;
    cronTimer(NDtimerCallback, meta.id, expiry);
}
function NDtimerCallback(docid) {
    delete dst_bucket[docid];
}
function timerCallback(docid) {
    dst_bucket[docid] = 'from timerCallback';
}
