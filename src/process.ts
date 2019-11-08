import * as vscode from 'vscode';
import * as path from 'path';
import {exec} from 'child_process';

export class Runner {
    private _config: Config;
    private _context: vscode.ExtensionContext;
    private _output: vscode.OutputChannel;

    constructor(context: vscode.ExtensionContext){
        this._context = context;
        this._config = <Config><any>vscode.workspace.getConfiguration("digitalonus.consulizer");
        this._output = vscode.window.createOutputChannel('digitalonus');
        
        this._output.append(JSON.stringify(this._config));
    }

}

interface Command {
    cmd: string;
    isAsync:boolean;
}

interface Config {
    clearConsole: boolean;
    commands: Array<Command>;
}
