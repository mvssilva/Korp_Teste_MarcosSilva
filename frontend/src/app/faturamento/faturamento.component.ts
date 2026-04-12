import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { FaturamentoService } from './faturamento.service';


@Component({
  selector: 'app-faturamento',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './faturamento.component.html',
  styleUrl: './faturamento.component.css'
})
export class FaturamentoComponent{

  itemAtualCodigo: string = '';
  itemAtualQuantidade: number = 1;
  itensNoCarrinho: any[] = [];

  mensagemFeedback: string = '';
  tipoFeedBack: 'sucesso' | 'erro' = 'sucesso';
  
  constructor(private faturamentoService: FaturamentoService){}
  
  adicionarItem(): void {
    if (this.itemAtualCodigo || this.itemAtualQuantidade > 0){
      this.itensNoCarrinho.push({
        codigo: this.itemAtualCodigo,
        itemAtualQuantidade: this.itemAtualQuantidade
      });
      this.itemAtualCodigo = '';
      this.itemAtualQuantidade = 1;
    }
  }

  salvarEImprimirNota() {
    if (this.itensNoCarrinho.length === 0) {
      this.mostrarMensagem('erro', 'Adicione itens antes de faturar.')
      return;
    }

    const payload = { itens: this.itensNoCarrinho };

    this.faturamentoService.emitirNota(payload).subscribe({
      next: (resposta) => {
        this.mostrarMensagem('sucesso', 'Nota fiscal emitida com sucesso!');
        this.itensNoCarrinho = [];
      },
      error: (erro) => {
        console.error('Erro ao faturar:', erro);
        this.mostrarMensagem('erro', 'Ocorreu um erro ao emitir a nota.')
      }
    });
  }

  mostrarMensagem(tipo: 'sucesso' | 'erro', texto: string) {
    this.tipoFeedBack = tipo;
    this.mensagemFeedback = texto;
    
    setTimeout(() => {
      this.mensagemFeedback = '';
    }, 3000);
  }

}
