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
    const produto = {
      codigo: this.novoCodigo,
      descricao: this.novaDescricao,
      saldo: this.novoSaldo
    }
    this.estoqueService.salvarProdutos(produto).subscribe(() => {
      this.carregarProdutos();
      this.novoCodigo = '';
      this.novaDescricao = '';
      this.novoSaldo = 0;
    })
  }
}