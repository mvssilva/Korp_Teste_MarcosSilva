import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common'; 
import { FormsModule } from '@angular/forms'; 
import { EstoqueService } from './estoque.service';

@Component({
  selector: 'app-estoque',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './estoque.component.html',
  styleUrl: './estoque.component.css'
})
export class EstoqueComponent implements OnInit{
  novoCodigo: string = '';
  novaDescricao: string = '';
  novoSaldo: number = 0;

  listaProdutos: any[] = [];

  // --- Controle do Toast de Notificação ---
  mensagemFeedback: string = '';
  tipoFeedBack: 'sucesso' | 'erro' = 'sucesso';

  constructor(private estoqueService: EstoqueService){}

  ngOnInit() {
    this.carregarProdutos();
  }

  carregarProdutos(){
    this.estoqueService.listarProdutos().subscribe(dados =>{
      this.listaProdutos = dados;
    });
  }

  salvar() {
    const codigoLimpo = this.novoCodigo.trim();
    const descricaoLimpa = this.novaDescricao.trim();

    if (!codigoLimpo || !descricaoLimpa) {
      this.mostrarMensagem('erro', 'Código e Descrição são obrigatórios!');
      return;
    }

    if (this.novoSaldo < 0) {
      this.mostrarMensagem('erro', 'O saldo não pode ser negativo.');
      return;
    }

    const produto = {
      codigo: codigoLimpo.toUpperCase(),
      descricao: descricaoLimpa,
      saldo: this.novoSaldo
    }

    this.estoqueService.salvarProdutos(produto).subscribe({
      next: () => {
        this.carregarProdutos();
        this.novoCodigo = '';
        this.novaDescricao = '';
        this.novoSaldo = 0;
        this.mostrarMensagem('sucesso', 'Produto cadastrado com sucesso!');
      },
      error: (erro) => {
        console.error('Erro ao salvar o produto:', erro);
        this.mostrarMensagem('erro', 'Falha ao salvar produto no banco.');
      }
    });
  }

  // --- Função que exibe a notificação na tela ---
  mostrarMensagem(tipo: 'sucesso' | 'erro', texto: string) {
    this.tipoFeedBack = tipo;
    this.mensagemFeedback = texto;
    
    // Oculta após 3 segundos
    setTimeout(() => {
      this.mensagemFeedback = '';
    }, 3000); 
  }
}