import { Component } from '@angular/core';
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
  // --- VARIÁVEIS DE CONTROLE DO PDF ---
  notaIniciada: boolean = false;
  notaFinalizada: boolean = false;
  numeroNotaExibicao: string = '---';
  statusNota: string = 'AGUARDANDO';
  
  // Controle de formulário
  itemAtualCodigo: string = '';
  itemAtualQuantidade: number = 1; 
  itensNoCarrinho: any[] = [];

  // Controle de UI
  mensagemFeedback: string = '';
  tipoFeedBack: 'sucesso' | 'erro' = 'sucesso';
  isProcessando: boolean = false;
  
  constructor(private faturamentoService: FaturamentoService) {}

  // Inicia o processo da nota
  iniciarNovaNota() {
    this.statusNota = 'CARREGANDO...'; // Mostra que está pensando
    this.numeroNotaExibicao = '---';

    this.faturamentoService.buscarProximoId().subscribe({
      next: (resposta) => {
        // Agora sim! O número real que veio do banco de dados
        this.numeroNotaExibicao = 'NF-' + resposta.proximo_id; 
        
        this.notaIniciada = true;
        this.notaFinalizada = false;
        this.statusNota = 'ABERTA';
        this.itensNoCarrinho = [];
        this.mostrarMensagem('sucesso', 'Nova Nota Fiscal aberta. Pode adicionar os itens.');
      },
      error: (erro) => {
        this.mostrarMensagem('erro', 'Falha ao comunicar com o servidor para buscar o número da nota.');
        this.statusNota = 'ERRO';
      }
    });
  }
  
  adicionarItem(): void {
    const codigoLimpo = this.itemAtualCodigo.trim().toUpperCase();

    if (codigoLimpo === '' || this.itemAtualQuantidade <= 0) {
      this.mostrarMensagem('erro', 'Preencha o código e a quantidade corretamente.');
      return;
    }

    this.faturamentoService.consultarProdutoNoEstoque(codigoLimpo).subscribe({
      next: (produto) => {
        this.itensNoCarrinho.push({
          produto_codigo: produto.codigo, 
          descricao: produto.descricao, 
          quantidade: this.itemAtualQuantidade
        });
        this.itemAtualCodigo = '';
        this.itemAtualQuantidade = 1;
      },
      error: (erro) => {
        this.mostrarMensagem('erro', `Produto ${codigoLimpo} não encontrado.`);
      }
    });
  }

    salvarEImprimirNota() {
    if (this.itensNoCarrinho.length === 0) return;

    this.isProcessando = true;
    let notaCriadaComSucesso = false; // Flag para rastrear o estado

    const payload = { 
      itens: this.itensNoCarrinho.map(item => ({
        produto_codigo: item.produto_codigo,
        quantidade: item.quantidade
      })) 
    };

    this.faturamentoService.criarNota(payload).pipe(
      switchMap((respostaCriacao: any) => {
        // Se chegamos aqui, a nota foi criada no banco!
        notaCriadaComSucesso = true; 
        const idDaNota = respostaCriacao.id || respostaCriacao.ID || respostaCriacao.numero;
        this.numeroNotaExibicao = 'NF-' + idDaNota;
        
        return this.faturamentoService.imprimirNota(idDaNota);
      })
    ).subscribe({
      next: (respostaFinal) => {
        this.mostrarMensagem('sucesso', 'Nota emitida com sucesso!');
        this.finalizarFluxoUI();
      },
      error: (erro) => {
        console.error('Erro:', erro);
        
        if (notaCriadaComSucesso) {
          // TRATAMENTO ESPECIAL: A nota existe, mas a impressão falhou.
          this.mostrarMensagem('erro', 'Nota gerada, mas houve erro na impressão/baixa. Verifique o estoque.');
          this.finalizarFluxoUI(); // Limpa o carrinho para evitar duplicidade!
        } else {
          // Erro na criação: O usuário pode tentar de novo.
          this.mostrarMensagem('erro', 'Erro ao criar a nota. Tente novamente.');
          this.isProcessando = false;
        }
      }
    });
  }

  // Função auxiliar para evitar repetição de código
  finalizarFluxoUI() {
    this.itensNoCarrinho = [];
    this.statusNota = 'FECHADA';
    this.notaFinalizada = true;
    this.isProcessando = false;
  }

  mostrarMensagem(tipo: 'sucesso' | 'erro', texto: string) {
    this.tipoFeedBack = tipo;
    this.mensagemFeedback = texto;
    setTimeout(() => {
      this.mensagemFeedback = '';
    }, 3000); 
  }
}