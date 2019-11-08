// The module 'vscode' contains the VS Code extensibility API
// Import the module and reference it with the alias vscode in your code below
import * as vscode from 'vscode';

// this method is called when your extension is activated
// your extension is activated the very first time the command is executed
export function activate(context: vscode.ExtensionContext) {

	console.log('Congratulations, your extension "consulizer" is now active!');
	let disposable = vscode.commands.registerCommand('extension.consulize', () => {
			var current = vscode.window.activeTextEditor!.document.uri.fsPath;
			if (current){
				vscode.window.showInformationMessage(current);
			}
	});

	context.subscriptions.push(disposable);
}

export function deactivate() {}
