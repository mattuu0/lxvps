function Init() {
    const decoder = new TextDecoder();

    // websocket に接続
    var ws = new WebSocket('/app/ws');

    // websocket からデータがきたら
    ws.onmessage = function (event) {
        // websocket からデータがきたら
        console.log(event.data);

        try {
            // でコードして書き込む
            term.write(decoder.decode(event.data));
        } catch (ex) {
            console.log(ex);
        }
    };


    // websocket を arraybuffer に変換
    ws.binaryType = 'arraybuffer';

    // ターミナルインスタンスを作成
    var term = new Terminal();

    // ターミナルを表示
    term.open(document.getElementById('terminal'));

    // プラグインをロード
    // プラグインを初期化
    const clipboardAddon = new ClipboardAddon.ClipboardAddon();
    const fitAddon = new FitAddon.FitAddon();

    // プラグインをロード
    term.loadAddon(clipboardAddon);
    term.loadAddon(fitAddon);

    // 画面サイズをフィットする
    fitAddon.fit();

    // ウィンドウがリサイズされたとき
    window.addEventListener('resize', function (evt) {
        // サイズを表示する
        console.log("resize", window.innerWidth, window.innerHeight);

        // フィットする
        fitAddon.fit();
    });

    // リサイズイベントをつける
    term.onResize(function (size) {
        console.log("new size", size);
        ws.send(JSON.stringify({ type: 'resize', "data": "", cols: size.cols, rows: size.rows }));
    });

    // キーが押された時のイベントを追加
    term.onKey(function (evt) {
        console.log(evt);
        // キーイベントを 送信する
        ws.send(JSON.stringify({ type: 'key', "data": evt.key, "cols": term.cols, "rows": term.rows }));
    });


    // 接続が開いたときに
    ws.onopen = function (event) {
        // キーイベントを 送信する
        ws.send(JSON.stringify({ type: 'resize', "data": "", "cols": term.cols, "rows": term.rows }));
    }
}

// 初期化
Init();