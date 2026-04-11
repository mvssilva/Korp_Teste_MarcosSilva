import { Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { EstoqueComponent } from './estoque/estoque.component';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [RouterOutlet, EstoqueComponent],
  templateUrl: './app.component.html', // CERTIFIQUE-SE QUE ESTÁ APP E NÃO ESTOQUE AQUI
  styleUrl: './app.component.css'
})
export class AppComponent {
  title = 'frontend';
}