import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { HttpClient } from '@angular/common/http';

@Component({
  selector: 'app-chat-ia',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './chat-ia.component.html',
  styleUrl: './chat-ia.component.css'
})
export class ChatIaComponent {
  chatAberto = false;
  mensagemUsuario = '';
  historico: { autor: 'user' | 'ia', texto: string }[] = [];
  isDigitando = false;

  private apiUrl = 'http://localhost:8081/chat';

  constructor(private http: HttpClient) {}

  toggleChat() {
    this.chatAberto = !this.chatAberto;
  }

  enviarMensagem() {
    if (!this.mensagemUsuario.trim() || this.isDigitando) return;

    // Adiciona a pergunta do usuário ao histórico local
    const pergunta = this.mensagemUsuario;
    this.historico.push({ autor: 'user', texto: pergunta });
    this.mensagemUsuario = '';
    this.isDigitando = true;

    // Consome a rota do Gemini que criamos no Go
    this.http.post<any>(this.apiUrl, { mensagem: pergunta }).subscribe({
      next: (res) => {
        this.historico.push({ autor: 'ia', texto: res.resposta });
        this.isDigitando = false;
      },
      error: () => {
        this.historico.push({ autor: 'ia', texto: 'Desculpe, estou com instabilidade agora.' });
        this.isDigitando = false;
      }
    });
  }
}