// The module 'vscode' contains the VS Code extensibility API
// Import the module and reference it with the alias vscode in your code below
import * as vscode from 'vscode';
import {Runner} from './process';

// this method is called when your extension is activated
// your extension is activated the very first time the command is executed
export function activate(context: vscode.ExtensionContext) {

	console.log('Congratulations, your extension "consulizer" is now active!');
	let disposable = vscode.commands.registerCommand('extension.consulize', () => {
			var current = vscode.window.activeTextEditor ? vscode.window.activeTextEditor.document.uri.fsPath : "";
			if (current){
				vscode.window.showInformationMessage(current);
				var runner = new Runner(context);
			}else{
				vscode.window.showInformationMessage('no working file');
			}
	});

	context.subscriptions.push(disposable);
}

export function deactivate() {}
