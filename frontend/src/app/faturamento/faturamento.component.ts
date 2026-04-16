import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { FaturamentoService } from './faturamento.service';
import { switchMap } from 'rxjs/operators';

@Component({
  selector: 'app-faturamento',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './faturamento.component.html',
  styleUrl: './faturamento.component.css'
})
export class FaturamentoComponent {

  // Controle de formulário
  itemAtualCodigo: string = '';
  itemAtualQuantidade: number = 1; 
  itensNoCarrinho: any[] = [];

  // Controle de UI (Notificações e Indicador de Processamento)
  mensagemFeedback: string = '';
  tipoFeedBack: 'sucesso' | 'erro' = 'sucesso';
  isProcessando: boolean = false;
  
  constructor(private faturamentoService: FaturamentoService) {}
  
  adicionarItem(): void {
    const codigoLimpo = this.itemAtualCodigo.trim().toUpperCase();

    if (codigoLimpo === '' || this.itemAtualQuantidade <= 0) {
      this.mostrarMensagem('erro', 'Preencha o código e a quantidade corretamente.');
      return;
    }

    // Consulta o microsserviço de Estoque antes de adicionar à tabela
    this.faturamentoService.consultarProdutoNoEstoque(codigoLimpo).subscribe({
      next: (produto) => {
        // Sucesso: O produto existe! Adiciona ao carrinho com a descrição.
        this.itensNoCarrinho.push({
          produto_codigo: produto.codigo, 
          descricao: produto.descricao, // Salva o nome para mostrar na tela
          quantidade: this.itemAtualQuantidade
        });
        
        // Limpa os campos
        this.itemAtualCodigo = '';
        this.itemAtualQuantidade = 1;
      },
      error: (erro) => {
        // Falha: O Estoque devolveu 404 Not Found
        this.mostrarMensagem('erro', `Produto ${codigoLimpo} não encontrado no estoque.`);
      }
    });
  }

  salvarEImprimirNota() {
    // Valida se o carrinho possui itens antes do envio
    if (this.itensNoCarrinho.length === 0) {
      this.mostrarMensagem('erro', 'Adicione itens antes de faturar.')
      return;
    }

    // Bloqueia múltiplas submissões e exibe o feedback visual no botão
    this.isProcessando = true;
    
    // Estrutura o payload (DTO) conforme o contrato esperado pela API
    const payload = { 
      itens: this.itensNoCarrinho.map(item => ({
        produto_codigo: item.produto_codigo,
        quantidade: item.quantidade
      })) 
    };

    // Fluxo RxJS: Cria a nota e, em seguida, dispara a impressão/baixa de estoque
    this.faturamentoService.criarNota(payload).pipe(
      switchMap((respostaCriacao: any) => {
        const idDaNota = respostaCriacao.id || respostaCriacao.ID || respostaCriacao.numero; 
        return this.faturamentoService.imprimirNota(idDaNota);
      })
    ).subscribe({
      next: (respostaFinal) => {
        // Sucesso em ambas as etapas (criação e impressão)
        this.mostrarMensagem('sucesso', 'Nota emitida e saldo abatido com sucesso!');
        this.itensNoCarrinho = []; 
        this.isProcessando = false; 
      },
      error: (erro) => {
        console.error('Erro no processo de faturamento:', erro);
        this.mostrarMensagem('erro', 'Ocorreu um erro ao processar a nota.');
        this.isProcessando = false; 
      }
    });
  }

  // Exibe notificação temporária (Toast) para o usuário
  mostrarMensagem(tipo: 'sucesso' | 'erro', texto: string) {
    this.tipoFeedBack = tipo;
    this.mensagemFeedback = texto;
    
    // Oculta a notificação após 3 segundos
    setTimeout(() => {
      this.mensagemFeedback = '';
    }, 3000); 
  }
}